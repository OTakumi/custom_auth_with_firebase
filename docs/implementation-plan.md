# Implementation Plan

This document outlines the phased approach for developing the Firebase custom authentication service.

## Phase 1: API Server Foundation [Completed]

1.  **Setup Project Structure:** Organized the Go server project with a standard directory structure (`cmd`, `internal`).
2.  **Create API Endpoint:** Implemented the initial API endpoint (`POST /auth/otp`) using the Gin framework to receive the user's email address.

## Phase 2: OTP Generation and Delivery [In Progress]

3.  **Implement OTP Generation Logic:** Implemented the functionality to generate a 6-digit one-time password. [Completed]
4.  **Cache OTP in Firestore:** Implemented the logic to store the generated OTP in Firebase Firestore with an expiration time. [Completed]
5.  **Write Integration Test:** Added an integration test to verify the OTP generation and Firestore storage logic using the Firebase Emulator. [Completed]
6.  **Implement Email Delivery:** Implement the functionality to send the OTP to the user's email address. This can be done using an external service like SendGrid.

## Phase 3: Firebase Integration and Authentication

7.  **Introduce Firebase Admin SDK:** Set up the Firebase Admin SDK in the Go project.
8.  **Implement OTP Verification and Custom Token Issuance:** Implement an API endpoint (`POST /auth/verify`) to verify the OTP entered by the user and issue a Firebase custom token upon success.

## Phase 4: Client-side Implementation

9.  **Create Frontend:** Implement a simple client with a UI for entering an email address and OTP.
10. **Implement API Integration and Firebase Sign-in:** Implement the logic on the client to call the API and use the received custom token to sign in with Firebase.

The development will continue with **Phase 2: Email Delivery**.
