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


## Successful Response Example

A successful request to the `/unwired/users/otp` endpoint will return a response like this:

**Status Code:** `201 Created`

**JSON Response:**

```json
{
  "status": "success",
  "data": {
    "one_time_passcode": "3JKLZV",
    "expires_at": 1751928133,
    "email": "test@nomail",
    "workvivo_user_id": 3983,
    "organisation_id": 165,
    "login_url": "https://unwired.workvivo.red/login"
  },
  "message": "One Time Passcode generated successfully."
}
```


## Rate Limit Response Example

If you exceed the allowed number of password reset requests, you will receive a response following `429 Too Many Requests`

For unwired.workvivo.red we have increased the rate limits:

* 10 requests over a 10 minute period
* 1440 requests over a 24 hour period


**Status Code:** `429 Too Many Requests`

**JSON Response:**

```json
{
  "status": "error",
  "data": null,
  "meta": {
    "errors": [
      {
        "path": "{}",
        "message": "You can only request 2 password resets in a 10-minute period."
      }
    ]
  }
}
```

**Status Code:** `429 Too Many Requests`

```json
{
  "status": "error",
  "data": null,
  "meta": {
    "errors": [
      {
        "path": "{}",
        "message": "You have reached the maximum number of password reset requests (10) for the past 24 hours."
      }
    ]
  }
}
```
