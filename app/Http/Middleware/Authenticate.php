<?php

namespace App\Http\Middleware;

use App\GlobalVar;

use Jose\Component\Core\AlgorithmManager;
use Jose\Component\Encryption\Algorithm\KeyEncryption\A128GCMKW;
use Jose\Component\Encryption\Algorithm\ContentEncryption\A128GCM;
use Jose\Component\Encryption\Compression\CompressionMethodManager;
use Jose\Component\Encryption\Compression\Deflate;
use Jose\Component\Encryption\JWEBuilder;
use Jose\Component\Encryption\JWEDecrypter;
use Jose\Component\Checker;
use Jose\Component\Checker\ClaimCheckerManager;
use Jose\Component\Checker\HeaderCheckerManager;
use Jose\Component\Checker\AlgorithmChecker;
use Jose\Component\Encryption\JWETokenSupport;

use Jose\Component\Checker\InvalidClaimException;
// use Jose\Component\Checker\InvalidHeaderException;

use Jose\Component\Core\JWK;
use Jose\Component\Encryption\Serializer\CompactSerializer;

use Illuminate\Http\Request;
use Closure;
use Exception;

class Authenticate
{
    public function handle(Request $request, Closure $next)
    {
        $shouldCheckRefreshToken = false;
        $shouldIssueNewTokenWithUID = null;

        $token = $request->headers->get('auth-token');
        $refresh_token = $request->headers->get('refresh-token');
        // 如果 短效token 不存在，返回
        if ($token == null) {
            return response()->json(GlobalVar::LACK_AUTH_ERR_RESPONSE);
        }

        $keyEncryptionAlgorithmManager = new AlgorithmManager([
            new A128GCMKW(),
        ]);
        $contentEncryptionAlgorithmManager = new AlgorithmManager([
            new A128GCM(),
        ]);
        $compressionMethodManager = new CompressionMethodManager([
            new Deflate(),
        ]);

        $jweDecrypter = new JWEDecrypter(
            $keyEncryptionAlgorithmManager,
            $contentEncryptionAlgorithmManager,
            $compressionMethodManager
        );

        $serializor = new CompactSerializer();

        $jwk = new JWK([
            'kty' => 'oct',
            'k' => env('APP_KEY'),
        ]);

        $headerCheckerManager = new HeaderCheckerManager(
            [
                new AlgorithmChecker(['A128GCMKW']),
            ],
            [
                new JWETokenSupport(),
            ]
        );

        $claimCheckerManager = new ClaimCheckerManager(
            [
                // new Checker\NotBeforeChecker(),
                new Checker\ExpirationTimeChecker()
            ]
        );


        // 检查 短效token
        try {
            $userJWE = $serializor->unserialize($token);

            // check方法在不成功的情况下是抛错而不是返回错误……这个东西看不懂
            $successheader = $headerCheckerManager->check($userJWE, 0, ['alg', 'enc']);
            $success = $jweDecrypter->decryptUsingKey($userJWE, $jwk, 0);

            $payload = json_decode($userJWE->getPayload(), true);
            $claimCheckerManager->check($payload);
        } catch (InvalidClaimException $err) {
            $shouldCheckRefreshToken = true;
        } catch (Exception $err) {
            // 如果 短效token 无效，返回
            return response()->json(GlobalVar::INVALID_AUTH_ERR_RESPONSE);
        }

        // 如果 短效token 过期，检查是否有 RefreshToken
        if ($shouldCheckRefreshToken) {
            if ($refresh_token === null) {
                return response()->json(GlobalVar::LACK_REFRESHTOKEN_ERR_RESPONSE);
            }


            try {
                $userRefreshJWE = $serializor->unserialize($refresh_token);
                $successheader = $headerCheckerManager->check($userRefreshJWE, 0, ['alg', 'enc']);
                $success = $jweDecrypter->decryptUsingKey($userRefreshJWE, $jwk, 0);

                $refreshTokenPayload = json_decode($userRefreshJWE->getPayload(), true);

                $claimCheckerManager->check($refreshTokenPayload);

                if (!array_key_exists('uid', $refreshTokenPayload)) {
                    // 等同 刷新token 无效
                    throw new Exception('无效refreshToken');
                }

                $shouldIssueNewTokenWithUID = $refreshTokenPayload['uid'];
            } catch (InvalidClaimException $err) {
                // 如果 刷新token 也过期，返回
                return response()->json(GlobalVar::EXPIRED_REFRESHTOKEN_ERR_RESPONSE);
            } catch (Exception $err) {
                // 如果 刷新token 也无效，返回
                return response()->json(GlobalVar::INVALID_REFRESHTOKEN_ERR_RESPONSE);
            }
        }

        // 转发到控制器处理逻辑
        $response = $next($request);

        // 检查是否需要重新签发
        if ($shouldIssueNewTokenWithUID === null) {
            return $response;
        }

        /**
         * 如果需要签发，则需要
         * 1. 重签普通token
         * 2. 检查 RefreshToken 是否在重签时间段内
         */
        $uid = $shouldIssueNewTokenWithUID;
        $time = time();

        $jweBuilder = new JWEBuilder(
            $keyEncryptionAlgorithmManager,
            $contentEncryptionAlgorithmManager,
            $compressionMethodManager
        );

        $jwePayload = json_encode([
            /** 20分钟 */
            'exp' => $time + 60 * 20,
            'uid' => $uid,
            ]);
        $refreshPayload = json_encode([
            /** 60天 */
            'exp' => $time,
            'uid' => $uid,
        ]);

        $jwe = $jweBuilder
            ->create()
            ->withPayload($payload)
            ->withSharedProtectedHeader([
                'alg' => 'A128GCMKW',
                'enc' => 'A128GCM',
                'zip' => 'DEF'
            ])
            ->addRecipient($jwk)
            ->build();


        return $response->header('token', $serializor->serialize($jwe, 0));
    }
}
