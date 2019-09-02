<?php

/**
 * 会根据request中带有的UID进行全新签发的中间件
 */

namespace App\Http\Middleware;

use App\GlobalVar;

use Jose\Component\Core\AlgorithmManager;
use Jose\Component\Encryption\Algorithm\KeyEncryption\A128GCMKW;
use Jose\Component\Encryption\Algorithm\ContentEncryption\A128GCM;
use Jose\Component\Encryption\Compression\CompressionMethodManager;
use Jose\Component\Encryption\Compression\Deflate;
use Jose\Component\Encryption\JWEBuilder;

use Jose\Component\Core\JWK;
use Jose\Component\Encryption\Serializer\CompactSerializer;

use Illuminate\Http\Request;
use Closure;
use Exception;

class AuthIssuer
{
    public function handle(Request $request, Closure $next)
    {
        $request->headers->remove('uid');

        $response = $next($request);

        $issueUID = $request->headers->get('uid');

        if (!$issueUID) {
            return $response;
        }

        $issueTokenWithUID = $request->headers->get('issueTokenWithUID');
        $issueRefreshTokenWithUID = $request->headers->get('issueRefreshTokenWithUID');

        if (!$issueTokenWithUID && !$issueRefreshTokenWithUID) {
            return $response;
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
            'exp' => $time + GlobalVar::ACCESS_TOKEN_EXP_TIME,
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

        $refreshPayload = json_encode([
            /** 60天 */
            'exp' => $time + GlobalVar::REFRESH_TOKEN_EXP_TIME,
            'uid' => $uid,
        ]);

        return $response->header('token', $serializor->serialize($jwe, 0));
    }
}
