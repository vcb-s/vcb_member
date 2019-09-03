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

        /**
         * 如果需要签发，则需要
         * 1. 重签普通token
         * 2. 检查 RefreshToken 是否在重签时间段内
         */
        $time = time();

        $jweBuilder = new JWEBuilder(
            $keyEncryptionAlgorithmManager,
            $contentEncryptionAlgorithmManager,
            $compressionMethodManager
        );

        $jwePayload = json_encode([
            'exp' => $time + GlobalVar::ACCESS_TOKEN_EXP_TIME,
            'uid' => $issueUID,
        ]);
        $token = $jweBuilder
            ->create()
            ->withPayload($jwePayload)
            ->withSharedProtectedHeader([
                'alg' => 'A128GCMKW',
                'enc' => 'A128GCM',
                'zip' => 'DEF'
            ])
            ->addRecipient($jwk)
            ->build();

        $response->headers->set('auth-token', $serializor->serialize($token, 0));

        $refreshPayload = json_encode([
            'exp' => $time + GlobalVar::REFRESH_TOKEN_EXP_TIME,
            'uid' => $issueUID,
        ]);
        $refreshToken = $jweBuilder
            ->create()
            ->withPayload($refreshPayload)
            ->withSharedProtectedHeader([
                'alg' => 'A128GCMKW',
                'enc' => 'A128GCM',
                'zip' => 'DEF'
            ])
            ->addRecipient($jwk)
            ->build();
        $response->headers->set('refresh-token', $serializor->serialize($refreshToken, 0));

        return $response;
    }
}
