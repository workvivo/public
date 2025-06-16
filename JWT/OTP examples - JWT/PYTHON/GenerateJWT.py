import os
import json
import time
import base64
import hashlib
import secrets
from pathlib import Path
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa
import jwt
import requests

# === CONFIGURATION ===
priv_key_file = Path(__file__).with_name('private.pem')
pub_key_file = Path(__file__).with_name('public.pem')

# === USER NEEDS TO CONFIGURE THESE SETTINGS HERE ===
org_id = 165
app_workvivo = '3981'
app_aud = 'unwired.workvivo.red'
user_email = 'test@nomail'
jwt_lifetime = 600
api_url = 'https://api-gateway.workvivo.red/v1/unwired/users/otp'
app_issuer = 'org1'
app_subject = 'app'

# === KEYPAIR GENERATION OR LOADING ===

#
#  You can generate your own RSA keypair manually using OpenSSL:
#  
#  # Generate a 4096-bit private key
#  openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096
#  
#  # Extract the public key from the private key
#  openssl rsa -in private.pem -pubout -out public.pem
#  
#  Place both files (private.pem and public.pem) in the same directory as this script.
#  The script will use them automatically instead of generating new keys.
#  


def generate_rsa_keypair():
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=4096
    )
    priv_pem = private_key.private_bytes(
        serialization.Encoding.PEM,
        serialization.PrivateFormat.PKCS8,
        serialization.NoEncryption()
    )
    pub_pem = private_key.public_key().public_bytes(
        serialization.Encoding.PEM,
        serialization.PublicFormat.SubjectPublicKeyInfo
    )

    priv_key_file.write_bytes(priv_pem)
    pub_key_file.write_bytes(pub_pem)

    print("Generated and saved new RSA keypair.")
    return private_key, priv_pem, pub_pem

def load_keys():
    priv_pem = priv_key_file.read_bytes()
    pub_pem = pub_key_file.read_bytes()
    private_key = serialization.load_pem_private_key(priv_pem, password=None)
    print("Loaded existing RSA keypair.")
    return private_key, priv_pem, pub_pem

if not priv_key_file.exists() or not pub_key_file.exists():
    private_key, private_pem, public_pem = generate_rsa_keypair()
else:
    private_key, private_pem, public_pem = load_keys()

# === base64url encode ===
def base64url_encode(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b'=').decode('utf-8')

# === BUILD JWKS ===
public_numbers = private_key.public_key().public_numbers()
e = base64url_encode(public_numbers.e.to_bytes((public_numbers.e.bit_length() + 7) // 8, byteorder='big'))
n = base64url_encode(public_numbers.n.to_bytes((public_numbers.n.bit_length() + 7) // 8, byteorder='big'))

jwk = {
    'kty': 'RSA',
    'alg': 'RS256',
    'use': 'sig',
    'n': n,
    'e': e,
}

thumbprint_input = json.dumps({'e': e, 'kty': 'RSA', 'n': n}, separators=(',', ':')).encode('utf-8')
thumbprint = base64url_encode(hashlib.sha256(thumbprint_input).digest())
jwk['kid'] = thumbprint
jwks = {'keys': [jwk]}

# === CREATE JWT ===
now = int(time.time())
payload = {
    'jti': secrets.token_hex(32),
    'iss': app_issuer,
    'sub': app_subject,
    'workvivo_id': app_workvivo,
    'aud': app_aud,
    'iat': now,
    'nbf': now,
    'exp': now + jwt_lifetime,
    'state': secrets.token_hex(32),
}

token = jwt.encode(
    payload,
    private_pem,
    algorithm='RS256',
    headers={'kid': jwk['kid']}
)

# === OUTPUT ===
print(f"\n=== JWT ===\n{token}")
print(f"\n=== kid ===\n{jwk['kid']}")
print(f"\n=== public.pem ===\n{public_pem.decode()}")
print(f"\n=== jwks ===\n{json.dumps(jwks, indent=4)}")

# prompt to confirm JWKS is uploaded or public key
input("\nPress [Enter] once you have updated your JWKS at:\nhttps://HOST/admin/developers/apps/manage\n")

# === SEND REQUEST ===
headers = {
    'Workvivo-Id': str(org_id),
    'x-workvivo-jwt': token,
    'x-workvivo-jwt-keyid': jwk['kid'],
    'Accept': 'application/json',
    'Content-Type': 'application/json',
}
post_body = json.dumps({'email': user_email})

response = requests.post(api_url, headers=headers, data=post_body, verify=False)
print(f"\n=== Response ({response.status_code}) ===\n{response.text}")
