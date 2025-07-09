# Keys Used for Sample Code

The public/private key pair provided here is **only for use with the code samples in this repository**. These keys are registered and associated with the domain `unwired.workvivo.red` for demonstration purposes.

> **Warning:**
> **Never use these keys in your own environments.** These keys are public and intended for testing only. Always generate and register a unique key pair for each Workvivo organisation.

---

## Generating a Public/Private Key Pair

You can use OpenSSL to generate key pairs. Below are typical commands to create a new key pair:

### 1. Generate a private key
This creates `private.pem`:
```sh
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096 > /dev/null 2>&1
```

### 2. Extract the public key
This creates `public.pem`:
```sh
openssl rsa -in private.pem -pubout -out public.pem
```

### 3. Create a JWKS file from your public key
Use the `create_jwks.py` script to generate the `jwks.json` file from your `public.pem` file:
```sh
python3 create_jwks.py > jwks.json
```

---

## Additional Notes

- The provided keys and scripts are for demonstration and testing only.
- **Never share your private keys.**
- For production, always use secure storage and handling for your private keys.
- Update your Workvivo organisation with your new JWKS as required.