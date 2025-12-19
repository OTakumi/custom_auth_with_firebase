# Implementation Plan

This document outlines the phased approach for developing the Firebase custom authentication service.

## Phase 1: API Server Foundation

1.  **Setup Project Structure:** Organize the Go server project with a standard directory structure (e.g., `cmd`, `internal`, `pkg`).
2.  **Create API Endpoint:** Implement the initial API endpoint (`POST /auth/otp`) to receive the user's email address.

## Phase 2: OTP Generation and Delivery

3.  **Implement OTP Generation Logic:** Implement the functionality to generate a 6-digit one-time password.
4.  **Cache OTP:** Implement the logic to store the generated OTP in a cache (e.g., Redis) with an expiration time.
5.  **Implement Email Delivery:** Implement the functionality to send the OTP to the user's email address. This can be done using an external service like SendGrid.

## Phase 3: Firebase Integration and Authentication

6.  **Introduce Firebase Admin SDK:** Set up the Firebase Admin SDK in the Go project.
7.  **Implement OTP Verification and Custom Token Issuance:** Implement an API endpoint (`POST /auth/verify`) to verify the OTP entered by the user and issue a Firebase custom token upon success.

## Phase 4: Client-side Implementation

8.  **Create Frontend:** Implement a simple client with a UI for entering an email address and OTP.
9.  **Implement API Integration and Firebase Sign-in:** Implement the logic on the client to call the API and use the received custom token to sign in with Firebase.

The development will start with **Phase 1: API Server Foundation**.
