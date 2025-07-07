# Workvivo Unwired Users One-Time Passcode Sample

This repository provides sample code for [Workvivo's Unwired Users One-Time Passcodes API](https://developer.workvivo.com/#aa34c835-aefb-4ff4-b1ad-232d00d37a9a).

## Quick Start

Start experimenting with the Unwired Users API using the sample code below.

- **No need to create a new key pair to run the sample code:**
  Keys are already provided and registered with the `unwired.workvivo.red` environment.
- **For your own Workvivo environments:**
  You must generate and register your own key pair.

> **Security Note:**
> - Never share your private keys.
> - Always generate a unique key pair for each Workvivo organisation.
> - Do NOT share keys across Workvivo environments or instances.

---

## Sample Code

- [Go Sample Code](GO/main.go)
- [PHP Sample Code](PHP/GenerateJWT.php)
- [Python Sample Code](PYTHON/GenerateJWT.py)

---

## Additional Notes

- Rate limits apply to the `/unwired/users/otp` endpoint.
- The provided scripts and code samples are for demonstration and testing purposes only.
- For production, always use secure storage and handling for your private keys.
- Update your Workvivo organisation with your new JWKS as required.