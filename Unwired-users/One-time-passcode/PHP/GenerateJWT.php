<?php

require __DIR__.'/vendor/autoload.php';
use Firebase\JWT\JWT;

/**
 * You will need to generate a public/private key pair for your application
 * the following Keys are should only be used for testing purposes and should not be used in production.
 */
$privKeyFile = __DIR__.'/../Keys/private.pem';
$pubKeyFile = __DIR__.'/../Keys/public.pem';
$jwksFile = __DIR__.'/../Keys/jwks.json';

/*
   The Organisation ID is the ID of the Workvivo organisation.
*/

$orgId = 165;

/*
    The App ID of the Workvivo app with unwiredotp.*.write permission.
*/

$appWorkvivo = '3981';

/*
 The Domain associated with the Workvivo app
*/

$appAud = 'unwired.workvivo.red';
/*
       This is the email address of the user you want a one-time passcode
*/

$userEmail = 'test@nomail';

/*
       The API Gateway URL for Workvivo, sample code is using api.workvivo.red, for production environments you will need to change this to the appropriate URL for your Workvivo environment.

       EU Production
       api.workvivo.com
       api.eu2.workvivo.com

       US Production
       api.workvivo.us
       api.us2.workvivo.us

       Middle East Production
       api.workvivo.me

    */

$apiUrl = 'https://api.workvivo.red/v1/unwired/users/otp'; //In case it is EU Production HOST should be api.workvivo.com

/*
        this will be the customer’s host typically, e.g. acme.com
    */

$appIssuer = 'org1';

/*
        this will always be 'app' as this API is always called as a Applicaton configured in Workvivo that the Partner Application is using to generate OTPs
    */
$appSubject = 'app';

$privatePem = file_get_contents($privKeyFile);
$publicPem = file_get_contents($pubKeyFile);
$jwks = file_get_contents($jwksFile);
$jwksData = json_decode(file_get_contents($jwksFile), true);

echo "Loaded existing RSA keypair.\n";

$kid = null;
if (isset($jwksData['keys'][0]['kid'])) {
    $kid = $jwksData['keys'][0]['kid'];
} else {
    echo "No kid found in JWKS.\n";
}

// === CREATE JWT ===
$now = time();
$payload = [
    'aud'          => $appAud,
    'workvivo_id'  => $appWorkvivo,
    'sub'          => $appSubject,
    'iss'          => $appIssuer,
    'nbf'          => $now,
    'iat'          => $now,
    'exp'          => $now + 60,
    'jti'          => bin2hex(random_bytes(32)),
    'state'        => bin2hex(random_bytes(32)),
];

$token = JWT::encode(
    $payload,
    $privatePem,
    'RS256',
    $kid
);

// === OUTPUT ===
echo "\n=== JWT Payload ===\n".json_encode($payload, JSON_PRETTY_PRINT)."\n";
echo "\n=== JWT encoded ===\n$token\n";
echo "\n=== kid ===\n$kid\n";
echo "\n=== public.pem ===\n$publicPem\n";
echo "\n=== jwks ===\n".json_encode($jwksData, JSON_PRETTY_PRINT)."\n";

// === SEND REQUEST ===

$headers = [
    "Workvivo-Id: $orgId",
    "x-workvivo-jwt: $token",
    "x-workvivo-jwt-keyid: {$kid}",
    'Accept: application/json',
    'Content-Type: application/json',
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
