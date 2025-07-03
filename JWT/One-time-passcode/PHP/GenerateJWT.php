<?php

require __DIR__ . '/vendor/autoload.php';
use Firebase\JWT\JWT;

// === CONFIGURATION ===

$privKeyFile = __DIR__ . '/private.pem';
$pubKeyFile  = __DIR__ . '/public.pem';

// === USER NEEDS TO CONFIGURE THESE SETTINGS HERE ===
$orgId       = 165; // This is the organisation ID
$appWorkvivo = '3981'; // APP ID from your Workvivo app with unwiredotp.*.write permission
$appAud      = 'unwired.workvivo.red';
$userEmail   = 'test@nomail';

$jwtLifetime = 600;
$apiUrl      = 'https://api-gateway.workvivo.red/v1/unwired/users/otp'; //In case it is EU Production HOST should be api.workvivo.com
$appIssuer   = 'org1';
$appSubject  = 'app';

// === KEYPAIR GENERATION OR LOADING ===

/*
 * You can generate your own RSA keypair manually using OpenSSL:
 *
 * # Generate a 4096-bit private key
 * openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096
 *
 * # Extract the public key from the private key
 * openssl rsa -in private.pem -pubout -out public.pem
 *
 * Place both files (private.pem and public.pem) in the same directory as this script.
 * The script will use them automatically instead of generating new keys.
 */

if (! file_exists($privKeyFile) || ! file_exists($pubKeyFile)) {
    // generate a new 4096-bit RSA keypair
    $res = openssl_pkey_new([
        'private_key_bits' => 4096,
        'private_key_type' => OPENSSL_KEYTYPE_RSA,
    ]);
    openssl_pkey_export($res, $privatePem);
    $keyDetails   = openssl_pkey_get_details($res);
    $publicPem    = $keyDetails['key'];

    // save to file
    file_put_contents($privKeyFile, $privatePem);
    file_put_contents($pubKeyFile,  $publicPem);

    echo "Generated and saved new RSA keypair.\n";
} else {
    $privatePem = file_get_contents($privKeyFile);
    $publicPem  = file_get_contents($pubKeyFile);
    $keyDetails = openssl_pkey_get_details(
        openssl_pkey_get_public($publicPem)
    );
    echo "Loaded existing RSA keypair.\n";
}

// === UTILITY: base64url ===

function base64url_encode(string $data): string {
    return rtrim(strtr(base64_encode($data), '+/', '-_'), '=');
}

// === BUILD JWKS ===

$rsa = $keyDetails['rsa'];
$jwk = [
    'kty' => 'RSA',
    'alg' => 'RS256',
    'use' => 'sig',
    'n'   => base64url_encode($rsa['n']),
    'e'   => base64url_encode($rsa['e']),
];
$thumbprint = base64url_encode(
    hash('sha256', json_encode([
        'e'   => $jwk['e'],
        'kty' => 'RSA',
        'n'   => $jwk['n'],
    ]), true)
);
$jwk['kid'] = $thumbprint;
$jwks = ['keys' => [$jwk]];

// === CREATE JWT ===

$now = time();
$payload = [
    'jti'          => bin2hex(random_bytes(32)),
    'iss'          => $appIssuer,
    'sub'          => $appSubject,
    'workvivo_id'  => $appWorkvivo,
    'aud'          => $appAud,
    'iat'          => $now,
    'nbf'          => $now,
    'exp'          => $now + $jwtLifetime,
    'state'        => bin2hex(random_bytes(32)),
];

$token = JWT::encode(
    $payload,
    $privatePem,
    'RS256',
    $jwk['kid']
);

// === OUTPUT ===

echo "\n=== JWT ===\n$token\n";
echo "\n=== kid ===\n{$jwk['kid']}\n";
echo "\n=== public.pem ===\n$publicPem\n";
echo "\n=== jwks ===\n" . json_encode($jwks, JSON_PRETTY_PRINT) . "\n";

// prompt to confirm JWKS is uploaded or public key
echo "\nPress [Enter] once you have updated your JWKS at:\n"
   . "https://HOST/admin/developers/apps/manage\n";
fgets(STDIN);

// === SEND REQUEST ===

$headers = [
    "Workvivo-Id: $orgId",
    "x-workvivo-jwt: $token",
    "x-workvivo-jwt-keyid: {$jwk['kid']}",
    "Accept: application/json",
    "Content-Type: application/json",
];
$postBody = json_encode(['email' => $userEmail]);

$ch = curl_init($apiUrl);
curl_setopt_array($ch, [
    CURLOPT_RETURNTRANSFER  => true,
    CURLOPT_POST            => true,
    CURLOPT_HTTPHEADER      => $headers,
    CURLOPT_POSTFIELDS      => $postBody,
    CURLOPT_SSL_VERIFYHOST  => false,
]);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

echo "\n=== Response ($httpCode) ===\n$response\n";
