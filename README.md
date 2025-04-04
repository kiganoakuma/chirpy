# Chirpy API Documentation

## Overview

Chirpy is a web application that provides a platform for users to post short messages called "chirps" (similar to tweets). It features user authentication, chirp management, and premium user upgrades via Polka payment processing.

## Technical Stack

- Backend: Go
- Database: PostgreSQL
- Authentication: JWT (JSON Web Tokens)
- Payment Processing: Polka API

## API Endpoints

### Health Check
- GET /api/healthz: Check if the API is running

### Chirps
- GET /api/chirps: Get all chirps
  - Optional query parameter: author_id to filter chirps by author
- GET /api/chirps/{chirpID}: Get a specific chirp by ID
- POST /api/chirps: Create a new chirp
- DELETE /api/chirps/{chirpID}: Delete a chirp

### Users
- POST /api/users: Create a new user
- PUT /api/users: Update user information

### Authentication
- POST /api/login: Authenticate a user and receive a JWT
- POST /api/refresh: Refresh an authentication token
- POST /api/revoke: Revoke an authentication token

### Payment Integration
- POST /api/polka/webhooks: Handle Polka payment webhooks for user upgrades

### Admin
- POST /admin/reset: Reset the application state (development only)
- GET /admin/metrics: Get application metrics

## Configuration

The application requires the following environment variables:
- DB_URL: PostgreSQL database connection string
- PLATFORM: Platform identifier
- JWT_SECRET: Secret key for JWT generation and validation
- POLKA_KEY: API key for Polka payment integration

## Running the Application

The application serves static files from the root directory at /app/ and listens on port 8080 by default.

# Start the server
go run main.go

## Security

- JWT-based authentication for protected endpoints
- Environment-based configuration for sensitive parameters
- Token refresh and revocation mechanisms

## Structure

The application follows a standard Go project layout:
- main.go: Entry point and server setup
- internal/database: Database interaction layer
- Various handler functions for API endpoints

Note: This documentation is a high-level overview and will be updated as more specific details about the application are provided..
