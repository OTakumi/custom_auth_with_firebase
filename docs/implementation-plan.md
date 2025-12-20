# Implementation Plan

This document outlines the phased approach for developing the Firebase custom authentication service.

## Phase 1: API Server Foundation [Completed]

1. **Setup Project Structure:** Organized the Go server project with a standard directory structure (`cmd`, `internal`).
2. **Create API Endpoint:** Implemented the initial API endpoint (`POST /auth/otp`) using the Gin framework to receive the user's email address.

## Phase 2: OTP Generation and Delivery [Completed]

3.  **Implement OTP Generation Logic:** Implemented the functionality to generate a 6-digit one-time password. [Completed]
4.  **Cache OTP in Firestore:** Implemented the logic to store the generated OTP in Firebase Firestore with an expiration time. [Completed]
5.  **Write Integration Test:** Added an integration test to verify the OTP generation and Firestore storage logic using the Firebase Emulator. [Completed]
6.  **Implement Email Delivery:** Implemented a dummy email delivery service that logs OTPs to the console. [Completed]

## Phase 3: Firebase Integration and Authentication [In Progress]

7.  **Setup Firebase Auth Client:** Initialize the Firebase Auth client in the `infrastructure/firebase` package to interact with the Firebase Authentication service. [Completed]
8.  **Implement OTP Verification Logic:** [Completed]
    -   **Repository:** Add a `Find` method to `OTPRepository` to retrieve OTP data from Firestore.
    -   **Usecase:** Add a `VerifyOTP` method to `OTPService` that validates the user's OTP against the stored one, checks for expiration, and deletes the OTP upon successful verification.
9.  **Implement Custom Token Generation:** [Completed]
    -   **Usecase:** Create a new `AuthService` responsible for generating a custom Firebase token for a given user ID (UID).
10. **Implement Verification Endpoint:** [Completed]
    -   **Handler:** Add a `VerifyOTP` method to `AuthHandler`.
    -   **Router:** Register a new `POST /auth/verify` route that accepts an email and OTP, and returns a custom Firebase token upon successful verification.
11. **Wire Dependencies:** Update `main.go` to initialize and inject the new services and repositories. [Completed]

## Phase 4: Client-side Implementation

12. **Create Frontend:** Implement a simple client with a UI for entering an email address and OTP.
13. **Implement API Integration and Firebase Sign-in:** Implement the logic on the client to call the API and use the received custom token to sign in with Firebase.

