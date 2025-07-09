import json
import base64
import hashlib
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend

def b64url_encode(data):
    return base64.urlsafe_b64encode(data).rstrip(b'=').decode('utf-8')

with open("public.pem", "rb") as f:
    pubkey_pem = f.read()
    pubkey = serialization.load_pem_public_key(pubkey_pem, backend=default_backend())

pubkey_der = pubkey.public_bytes(
    encoding=serialization.Encoding.DER,
    format=serialization.PublicFormat.SubjectPublicKeyInfo
)
kid = base64.urlsafe_b64encode(hashlib.sha256(pubkey_der).digest()).rstrip(b'=').decode('utf-8')

numbers = pubkey.public_numbers()
e = numbers.e.to_bytes((numbers.e.bit_length() + 7) // 8, byteorder="big")
n = numbers.n.to_bytes((numbers.n.bit_length() + 7) // 8, byteorder="big")

jwk = {
    "kty": "RSA",
    "use": "sig",
    "alg": "RS256",
    "kid": kid,
    "n": b64url_encode(n),
    "e": b64url_encode(e),
}

jwks = {"keys": [jwk]}

print(json.dumps(jwks, indent=2))