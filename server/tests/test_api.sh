#!/bin/bash

API_URL="http://localhost:8000"

echo "Testing /health endpoint..."
curl -v ${API_URL}/health

echo -e "\nTesting /auth/otp endpoint with a valid email..."
curl -v -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com"}' ${API_URL}/auth/otp

echo -e "\nTesting /auth/otp endpoint with an invalid email format (API will validate later)..."
curl -v -X POST -H "Content-Type: application/json" -d '{"email": "invalid-email"}' ${API_URL}/auth/otp

echo -e "\nTesting /auth/otp endpoint with missing email in body..."
curl -v -X POST -H "Content-Type: application/json" -d '{}' ${API_URL}/auth/otp

echo -e "\nTesting /auth/otp endpoint with invalid JSON..."
curl -v -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com",}' ${API_URL}/auth/otp

echo -e "\nTesting /auth/otp endpoint with GET method (should fail)..."
curl -v -X GET ${API_URL}/auth/otp
