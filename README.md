# Central Authorization RBAC Server

A high-performance, secure Central Authorization and Role-Based Access Control (RBAC) server built in Go. This service provides a centralized hub for managing authentication, users, roles, and permissions across an entire ecosystem of microservices.

It issues strictly standardized JSON Web Tokens (JWTs) secured with **RS256 Asymmetric Encryption**. 

## Features

- **Microservice Integration**: Generate `X-API-Key` values for different projects to isolate user bases and enforce strict access control.
- **Asymmetric Encryption (RS256)**: Central Auth signs tokens with a private RSA key. External microservices fetch the public key to verify tokens independently, eliminating the risk of symmetric key leakage.
- **RBAC**: Define fine-grained Roles and Permissions.

## How to Integrate With Your Projects

If you are building a new microservice (e.g., Python, Node.js, Go) and want to use this Central Auth server for user logins, follow these exact steps:

### 1. Register Your Project
To secure communication between your microservice and the Central Auth server, you must first create a Project and obtain an API Key.

Make a `POST` request to the Central Auth Server to register your project:
```bash
curl -X POST http://localhost:9000/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "My New App", "description": "Microservice for data processing"}'
```
*Note the returned `project_id`.*

Then, generate an API Key for that project:
```bash
curl -X POST http://localhost:9000/api/v1/projects/{project_id}/api-key
```
**Store this API Key in your microservice's `.env` file as `CENTRAL_AUTH_API_KEY`.**

### 2. Forward Registration & Login Requests
Your microservice should provide its own UI/endpoints for Login and Registration. When a user submits their credentials to your microservice, your microservice should forward those credentials to the Central Auth server.

**IMPORTANT:** You must include the `X-API-Key` header in all requests.

#### Register a User
```http
POST /api/v1/auth/register
Host: central-auth-server.com
X-API-Key: pk_live_your_api_key_here
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Login a User
```http
POST /api/v1/auth/login
Host: central-auth-server.com
X-API-Key: pk_live_your_api_key_here
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```
**Response:** You will receive a `token` (a JWT). Pass this token back to your user's browser (e.g., as an HTTP-only cookie or Bearer token).

### 3. Verify JWTs in Your Microservice (RS256)
When your user makes a protected request to your microservice, they will send the JWT. Your microservice must verify that this token is valid and was genuinely signed by the Central Auth Server.

Because we use RS256, your microservice needs the **Public Key** from Central Auth.

1. **Fetch the Public Key:**
   Your microservice should make a simple GET request to fetch the public key (usually on startup or lazily cached on the first request):
   ```http
   GET /api/v1/auth/public-key
   Host: central-auth-server.com
   ```
   *This endpoint is unprotected and does not require an API key.*

2. **Decode and Verify:**
   Use your language's standard JWT library (like `PyJWT` in Python or `jsonwebtoken` in Node) to decode the token using the fetched Public Key, explicitly enforcing the `RS256` algorithm.

#### Python Example
```python
import jwt
import httpx

# Fetch public key once and cache it
public_key = httpx.get("http://central-auth-server.com/api/v1/auth/public-key").text

# Verify token from user request
payload = jwt.decode(
    user_provided_token, 
    public_key, 
    algorithms=["RS256"]
)
print("User ID:", payload["sub"])
print("Email:", payload["email"])
```

## Running the Server Locally

```bash
# Install dependencies
go mod download

# Start the server (will automatically generate private_key.pem and public_key.pem if they don't exist)
go run main.go
```

## Running with Docker (Production)

```bash
docker compose up -d --build
```
