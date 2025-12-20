# End-to-End (E2E) Test Plan

This document outlines the test scenarios for the API endpoints of the Firebase custom authentication service.

## Prerequisites

- The backend service and its dependencies (Firebase emulator) are running, orchestrated via `docker compose`.
- The API is accessible at `http://localhost:8000`.
- The health check endpoint `GET /health` returns a `200 OK` status.

## Test Scenarios

### Endpoint: `POST /auth/otp`

This endpoint is responsible for generating and sending a one-time password.

| Test Case Name                     | Type      | Verification Steps                                                                                               | Expected Result                                                                                                    |
| ---------------------------------- | --------- | ---------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| **Request OTP with Valid Email**   | Normal    | 1. Send a `POST` request to `/auth/otp` with a valid JSON body, e.g., <br> `curl -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com"}' http://localhost:8000/auth/otp` | - **Status Code:** `200 OK`<br>- **Response Body:** Contains a success message and the generated OTP string.        |
| **Request OTP with Invalid Email** | Abnormal  | 1. Send a `POST` request to `/auth/otp` with a JSON body containing an invalid email, e.g., `{"email": "test.com"}`. | - **Status Code:** `400 Bad Request`<br>- **Response Body:** Contains an "Invalid email format" error message.      |
| **Request OTP with Missing Email** | Abnormal  | 1. Send a `POST` request to `/auth/otp` with an empty JSON body `{}`.                                            | - **Status Code:** `400 Bad Request`<br>- **Response Body:** Contains an "Invalid request body" error message.     |
| **Request OTP with Bad JSON**      | Abnormal  | 1. Send a `POST` request to `/auth/otp` with a malformed JSON string.                                            | - **Status Code:** `400 Bad Request`<br>- **Response Body:** Contains an "Invalid request body" error message.     |

---

### Endpoint: `POST /auth/verify`

This endpoint verifies the OTP and returns a custom Firebase token.

| Test Case Name                            | Type      | Verification Steps                                                                                                                                                                                                                                                               | Expected Result                                                                                                                 |
| ----------------------------------------- | --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **Verify with Correct OTP**               | Normal    | 1. Request an OTP for a valid email (e.g., `test@example.com`) by calling `/auth/otp`: <br> `curl -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com"}' http://localhost:8000/auth/otp` <br> (Note: You will need to manually extract the OTP from the response.)<br>2. Send a `POST` request to `/auth/verify` with the same email and the extracted OTP: <br> `curl -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com", "otp": "<EXTRACTED_OTP>"}' http://localhost:8000/auth/verify` | - **Status Code:** `200 OK`<br>- **Response Body:** Contains a `token` field with a valid Firebase custom token (JWT string).    |
| **Verify with Incorrect OTP**             | Abnormal  | 1. Request an OTP for a valid email by calling `/auth/otp`.<br>2. Send a `POST` request to `/auth/verify` with the same email but an incorrect OTP (e.g., `"000000"`).                                                                                                                     | - **Status Code:** `401 Unauthorized`<br>- **Response Body:** Contains an "Invalid or expired OTP" error message.               |
| **Verify with Expired OTP**               | Abnormal  | 1. Request an OTP for a valid email by calling `/auth/otp`.<br>2. Wait for 5 minutes for the OTP to expire.<br>3. Send a `POST` request to `/auth/verify` with the same email and the now-expired OTP.                                                                                     | - **Status Code:** `401 Unauthorized`<br>- **Response Body:** Contains an "Invalid or expired OTP" error message.               |
| **Verify without Requesting OTP**         | Abnormal  | 1. Send a `POST` request to `/auth/verify` with a valid email and a random OTP without having called `/auth/otp` first for that email.                                                                                                                                                 | - **Status Code:** `401 Unauthorized`<br>- **Response Body:** Contains an "Invalid or expired OTP" error message.               |
| **Verify Using the Same OTP Twice**       | Abnormal  | 1. Request an OTP for a valid email by calling `/auth/otp`.<br>2. Verify successfully once using the correct OTP.<br>3. Attempt to verify a second time using the exact same OTP.                                                                                                     | - **Status Code:** `401 Unauthorized`<br>- **Response Body:** Contains an "Invalid or expired OTP" error message.               |
| **Verify with Missing `otp` field**       | Abnormal  | 1. Send a `POST` request to `/auth/verify` with a JSON body containing only the `email` field.                                                                                                                                                                                     | - **Status Code:** `400 Bad Request`<br>- **Response Body:** Contains an "Invalid request body" error message.                  |
| **Verify with Missing `email` field**     | Abnormal  | 1. Send a `POST` request to `/auth/verify` with a JSON body containing only the `otp` field.                                                                                                                                                                                      | - **Status Code:** `400 Bad Request`<br>- **Response Body:** Contains an "Invalid request body" error message.                  |
