# Chirpy

Chirpy is a backend REST API and server platform for a social network inspired by X (formerly Twitter). It allows users to securely sign up, log in, post chirps, and view posts using validated JSON Web Tokens (JWT) and refresh tokens.

## Features

### Users

- Create an account.
- Log in to your account.
- Update your email address and change your password.
- Subscribe to Chirpy Red (premium account status via webhook integration).
- Post chirps.
- Read all chirps or filter chirps by a specific user.
- Sort chirps by post date and time in ascending or descending order (`?sort=asc` or `?sort=desc`).
- Delete your own chirps.

### Admin & System

Restricted administrative and system endpoints.

- Check server health and readiness status.
- View homepage access metrics.
- Reset all database data (development environment only).
- Revoke refresh tokens.

## Tech Stack & Requirements

- **Language:** Go (1.22+ recommended)
- **Database:** PostgreSQL
- **Database Migrations:** Goose
- **SQL Generation & Queries:** SQLC
- **Environment Variable Loading:** Godotenv

## Installation

### 1. Clone the Repository

Clone or download the project repository from GitHub, then navigate to the project root directory.

```bash
git clone <repository-url>
cd chirpy
```

Replace `<repository-url>` with the URL of your repository.

### 2. Install PostgreSQL

**macOS**

```bash
brew install postgresql@15
```

**Linux / WSL (Debian/Ubuntu)**

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

### 3. Install Goose

Goose is used to manage database migrations.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. Install SQLC

SQLC generates type-safe Go code from SQL queries.

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 5. Install Godotenv 

Godotenv loads environment variables from a local `.env` file during development.

```bash
go get github.com/joho/godotenv
```

## Configuration

1. Create a `.env` file in the project root directory.
2. Ensure `.env` is listed in your `.gitignore` file to prevent committing secrets, database credentials, or API keys.
3. Add the following environment variables to your `.env` file:

```dotenv
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="a_nice_long_random_string"
POLKA_KEY="YOUR_POLKA_KEY"
```

### Environment Variables

| Variable | Description |
| --- | --- |
| `DB_URL` | PostgreSQL connection string. |
| `PLATFORM` | Application environment. Set to `dev` for development. |
| `SECRET` | Secret used to sign and validate JWTs. |
| `POLKA_KEY` | API key used to authenticate Polka webhook requests. |

**Important:**
- Replace `username` and `password` with your PostgreSQL credentials.
- Replace `YOUR_POLKA_KEY` with your Polka API key.
- Generate a secure secret using the following command:

```bash
openssl rand -base64 64
```

Never commit your `.env` file or expose your application secrets.

## Database Setup

### 1. Start PostgreSQL

**macOS**

```bash
brew services start postgresql@15
```

**Linux**

```bash
sudo service postgresql start
```

### 2. Create the Database

Create the `chirpy` database if it does not already exist.

```bash
createdb chirpy
```

### 3. Run Database Migrations

Use Goose to apply the database migrations.

```bash
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/chirpy?sslmode=disable" up
```

Ensure the connection string matches the `DB_URL` value in your `.env` file and that `sql/schema` points to your project's migration directory.

## Running the Application

### Run Directly

Start the application using Go:

```bash
go run .
```

### Build and Run

Alternatively, build and run the application binary:

```bash
go build -o chirpy
./chirpy
```

The examples below assume the server is running at `http://localhost:8080`.

### macOS Binary Security Note

If you downloaded a precompiled macOS binary and receive a security warning such as *"Apple could not verify..."*, follow these steps:

1. Click **Done** or **Cancel** on the alert. Do not move the application to the Trash.
2. Open **System Settings → Privacy & Security**.
3. Scroll to the Security section and locate the blocked application message.
4. Click **Open Anyway**.
5. Authenticate using your Mac administrator password or Touch ID if prompted.

## API Reference

### 1. Health Check

Checks whether the server is running.

| Property | Value |
| --- | --- |
| **Method** | `GET` |
| **Endpoint** | `/api/healthz` |
| **Authentication** | None |
| **Success Response** | `200 OK` |
| **Common Errors** | `500 Internal Server Error` |

**Example Request**

```bash
curl http://localhost:8080/api/healthz
```

---

### 2. Create User

Registers a new user account.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/users` |
| **Authentication** | None |
| **Success Response** | `201 Created` |
| **Common Errors** | `400 Bad Request` |

**Request Body**

```json
{
  "email": "user@example.com",
  "password": "secure-password"
}
```

**Example Request**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secure-password"}'
```

---

### 3. Log In

Authenticates a user and returns their account details, access token, and refresh token.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/login` |
| **Authentication** | None |
| **Success Response** | `200 OK` |
| **Common Errors** | `401 Unauthorized` |

**Request Body**

```json
{
  "email": "user@example.com",
  "password": "secure-password"
}
```

**Example Request**

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secure-password"}'
```

---

### 4. Refresh Access Token

Generates a new access token using a valid refresh token.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/refresh` |
| **Authentication** | Bearer refresh token |
| **Success Response** | `200 OK` |
| **Common Errors** | `401 Unauthorized` |

**Example Request**

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN"
```

---

### 5. Revoke Refresh Token

Revokes a refresh token, preventing it from being used to obtain new access tokens.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/revoke` |
| **Authentication** | Bearer refresh token |
| **Success Response** | `204 No Content` |
| **Common Errors** | `401 Unauthorized` |

**Example Request**

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN"
```

---

### 6. Update User Account

Updates the authenticated user's email address and password.

| Property | Value |
| --- | --- |
| **Method** | `PUT` |
| **Endpoint** | `/api/users` |
| **Authentication** | Bearer access token |
| **Success Response** | `200 OK` |
| **Common Errors** | `400 Bad Request`, `401 Unauthorized` |

**Request Body**

```json
{
  "email": "new@example.com",
  "password": "new-secure-password"
}
```

**Example Request**

```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com","password":"new-secure-password"}'
```

---

### 7. Create a Chirp

Creates a new chirp for the authenticated user.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/chirps` |
| **Authentication** | Bearer access token |
| **Success Response** | `201 Created` |
| **Common Errors** | `400 Bad Request`, `401 Unauthorized` |

**Request Body**

```json
{
  "body": "A message from Chirpy"
}
```

**Example Request**

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"A message from Chirpy"}'
```

---

### 8. List Chirps

Retrieves all chirps, optionally filtered by author and sorted by creation time.

| Property | Value |
| --- | --- |
| **Method** | `GET` |
| **Endpoint** | `/api/chirps` |
| **Authentication** | None |
| **Success Response** | `200 OK` |
| **Common Errors** | `400 Bad Request` |

**Query Parameters**

| Parameter | Description | Accepted Values |
| --- | --- | --- |
| `author_id` | Filters chirps by the author's UUID. | Valid user UUID |
| `sort` | Sorts chirps by creation time. | `asc`, `desc` |

**Example Requests**

Retrieve all chirps:

```bash
curl http://localhost:8080/api/chirps
```

Retrieve chirps by a specific author:

```bash
curl "http://localhost:8080/api/chirps?author_id=AUTHOR_UUID"
```

Retrieve chirps by a specific author in descending order:

```bash
curl "http://localhost:8080/api/chirps?author_id=AUTHOR_UUID&sort=desc"
```

---

### 9. Get a Single Chirp

Retrieves a specific chirp by its UUID.

| Property | Value |
| --- | --- |
| **Method** | `GET` |
| **Endpoint** | `/api/chirps/{chirpID}` |
| **Authentication** | None |
| **Success Response** | `200 OK` |
| **Common Errors** | `404 Not Found` |

**Example Request**

```bash
curl http://localhost:8080/api/chirps/CHIRP_UUID
```

---

### 10. Delete a Chirp

Deletes a chirp belonging to the authenticated user.

| Property | Value |
| --- | --- |
| **Method** | `DELETE` |
| **Endpoint** | `/api/chirps/{chirpID}` |
| **Authentication** | Bearer access token |
| **Authorization** | The authenticated user must be the chirp's author. |
| **Success Response** | `204 No Content` |
| **Common Errors** | `401 Unauthorized`, `403 Forbidden`, `404 Not Found` |

**Example Request**

```bash
curl -X DELETE http://localhost:8080/api/chirps/CHIRP_UUID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 11. Polka Webhook

Upgrades a user's account to Chirpy Red when a valid Polka webhook event is received.

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/api/polka/webhooks` |
| **Authentication** | API key |
| **Success Response** | `204 No Content` |
| **Common Errors** | `401 Unauthorized`, `404 Not Found` |

**Request Body**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "USER_UUID"
  }
}
```

**Example Request**

```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Authorization: ApiKey YOUR_POLKA_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "event": "user.upgraded",
    "data": {
      "user_id": "USER_UUID"
    }
  }'
```

---

### 12. Admin: View Homepage Metrics

Displays homepage access metrics, including the number of times the homepage has been accessed.

| Property | Value |
| --- | --- |
| **Method** | `GET` |
| **Endpoint** | `/admin/metrics` |
| **Authentication** | None |
| **Success Response** | `200 OK` (HTML page) |
| **Common Errors** | `500 Internal Server Error` |

**Example Request**

```bash
curl http://localhost:8080/admin/metrics
```

---

### 13. Admin: Reset Database

Resets the database. This endpoint is available only when the application is running in the development environment (`PLATFORM=dev`).

| Property | Value |
| --- | --- |
| **Method** | `POST` |
| **Endpoint** | `/admin/reset` |
| **Authentication** | None |
| **Success Response** | `200 OK` |
| **Common Errors** | `403 Forbidden` |

**Example Request**

```bash
curl -X POST http://localhost:8080/admin/reset
```

> **Warning:** This endpoint resets database data. Use it only in a development environment.

## Example Workflow

The following workflow demonstrates a typical user journey, from account creation to authentication and posting chirps.

### 1. Create an Account

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secure-password"}'
```

### 2. Log In

Save the `token` and `refresh_token` returned in the response.

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secure-password"}'
```

### 3. Create a Chirp

Replace `ACCESS_TOKEN` with the access token returned during login.

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello, Chirpy!"}'
```

### 4. Fetch All Chirps

```bash
curl http://localhost:8080/api/chirps
```

### 5. Fetch Chirps by User

Replace `USER_UUID` with the UUID of the user whose chirps you want to retrieve.

```bash
curl "http://localhost:8080/api/chirps?author_id=USER_UUID"
```

### 6. Refresh Your Access Token

Use your refresh token to obtain a new access token.

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer REFRESH_TOKEN"
```

### 7. Log Out

Revoke your refresh token to prevent it from being used to obtain new access tokens.

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer REFRESH_TOKEN"
```
