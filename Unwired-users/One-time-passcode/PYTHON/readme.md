# Python Sample: Generate JWT for Workvivo Unwired Users

This script demonstrates how to generate a JWT and interact with the Workvivo Unwired Users API using Python.

## Setup & Usage

1. **Create and activate a virtual environment:**
   ```sh
   python3 -m venv venv
   source venv/bin/activate
   ```

2. **Install required dependencies:**
   ```sh
   pip install -r requirements.txt
   ```

3. **Run the script:**
   ```sh
   python3 GenerateJWT.py
   ```

## Notes
- Make sure `private.pem` and `public.pem` are present in the `../Keys/` directory relative to this script.
- The script will output the JWT, Key ID, public key, JWKS, and the API response.
- For security, never share your private key and always use a unique key pair for each environment.

