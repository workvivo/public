# Workvivo Unwired Users One-Time Passcode Sample

This repository provides sample code for [Workvivo's Unwired Users One-Time Passcodes API](https://developer.workvivo.com/#aa34c835-aefb-4ff4-b1ad-232d00d37a9a).

---

## Table of Contents
- [Quick Start](#quick-start)
- [Sample Code](#sample-code)
- [Additional Notes](#additional-notes)
- [Successful Response Example](#successful-response-example)
- [Rate Limit Information & Responses](#rate-limit-information--responses)

---

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

---

## Successful Response Example

A successful request to the `/unwired/users/otp` endpoint will return a response like this:

**Status Code:** `201 Created`

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

---

## Rate Limit Information & Responses

If you exceed the allowed number of password reset requests, you will receive a `429 Too Many Requests` response.

### Rate Limits by Environment

| Environment                | 10 min window | 24 hour window |
|----------------------------|:-------------:|:--------------:|
| **unwired.workvivo.red**   | 10 requests   | 1440 requests  |
| **Default (all others)**   | 2 requests    | 10 requests    |

> **Note:** The above limits apply to the `/unwired/users/otp` endpoint. Limits may be subject to change by Workvivo.

### Example Rate Limit Responses

**Exceeded 10-minute window (default):**

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

**Exceeded 24-hour window (default):**

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

---

## Troubleshooting

- Ensure your keys are valid and registered with your Workvivo environment.
- If you see SSL warnings in Python, see the comments in the sample code for how to disable them (not recommended for production).
- For any issues with dependencies, see the language-specific README files in each sample directory.

---

For more information, visit the [Workvivo Developer Documentation](https://developer.workvivo.com/).
