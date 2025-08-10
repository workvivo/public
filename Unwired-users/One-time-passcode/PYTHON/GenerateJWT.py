import json
import time
import base64
import secrets
from pathlib import Path
from cryptography.hazmat.primitives import serialization
import jwt
import requests
import urllib3

""""
 Disable SSL warnings for insecure requests (not recommended for production)
 This is only for development purposes or you have TLS Inspection enabled via Proxy. 
 Note: This is not recommended for production environments as it disables SSL verification.
"""
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

priv_key_file = Path(__file__).parent.parent / 'Keys' / 'private.pem'
pub_key_file = Path(__file__).parent.parent / 'Keys' / 'public.pem'
jwks_file = Path(__file__).parent.parent / 'Keys' / 'jwks.json'

""""
 The Organisation ID is the ID of the Workvivo organisation.
"""
org_id = 165

""""
 The App ID of the Workvivo app with unwiredotp.*.write permission.
"""
app_workvivo = '3981'

""""
 The Domain associated with the Workvivo app
"""

app_aud = 'unwired.workvivo.red'

""""
 	   EU Production
	   api.workvivo.com
	   api.eu2.workvivo.com

	   US Production
	   api.workvivo.us
	   api.us2.workvivo.us

	   Middle East Production
	   api.workvivo.me
"""

api_url = 'https://api.workvivo.red/v1/unwired/users/otp'

"""
 This is the email address of the user you want a one-time passcode
"""
user_email = 'test@nomail'

"""
 This will be the customer’s host typically, e.g. acme.com
"""
app_issuer = 'org1'

"""
 This will always be 'app' as this API is always called as an Application configured in Workvivo that the Partner Application is using to generate OTPs
"""
app_subject = 'app'

"""
 JWT lifetime in seconds
"""

jwt_lifetime = 600

"""
Load existing keys
"""
def load_keys():
    priv_pem = priv_key_file.read_bytes()
    pub_pem = pub_key_file.read_bytes()
    private_key = serialization.load_pem_private_key(priv_pem, password=None)
    return private_key, priv_pem, pub_pem

private_key, private_pem, public_pem = load_keys()

# === base64url encode ===
def base64url_encode(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b'=').decode('utf-8')

# === Read JWKS and extract kid ===
def read_jwks_and_get_kid():
    if not jwks_file.exists():
        return None
    with open(jwks_file, 'r') as f:
        jwks = json.load(f)
    keys = jwks.get('keys', [])
    if keys and 'kid' in keys[0]:
        return keys[0]['kid']
    return None

def main():
    # Read kid from JWKS file
    with open(jwks_file, 'r') as f:
        jwks = json.load(f)
    keys = jwks.get('keys', [])
    kid = keys[0]['kid'] if keys and 'kid' in keys[0] else None


      # === CREATE JWT ===
    now = int(time.time())
    payload = {
        'aud': app_aud,
        'app_id': app_workvivo,
        'sub': app_subject,
        'iss': app_issuer,
        'nbf': now,
        'iat': now,
        'exp': now + jwt_lifetime,
        'jti': secrets.token_hex(32),
        'state': secrets.token_hex(32),
    }
    token = jwt.encode(
        payload,
        private_pem,
        algorithm='RS256',
        headers={'kid': kid}
    )

    print(f"\nJWT payload:")
    print(json.dumps(payload, indent=2))
    print(f"\nJWT encoded:\n{token}\n")
    print(f"KeyID:\n{kid}\n")
    print("Public Key (PEM):")
    print(public_pem.decode())
    print("\nJWKS:")
    print(json.dumps(jwks, indent=2))

    # === SEND REQUEST ===
    headers = {
        'Workvivo-Id': str(org_id),
        'x-workvivo-jwt': token,
        'x-workvivo-jwt-keyid': kid,
        'Accept': 'application/json',
        'Content-Type': 'application/json',
    }
    post_body = json.dumps({'email': user_email})

    response = requests.post(api_url, headers=headers, data=post_body, verify=False)
    print(f"\nResponse:\n{response.text}")
    print(f"HTTP Code:\n{response.status_code}")

if __name__ == "__main__":
    main()
