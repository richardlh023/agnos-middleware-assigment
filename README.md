# Agnos Middleware System

Hospital Middleware System built with Go, Gin, and PostgreSQL. This system provides APIs for hospital staff authentication and patient search functionality with integration to external Hospital Information System (HIS) APIs.

## What This Project Does

This middleware system enables:

- **Staff Management**: Create and authenticate hospital staff accounts
- **Patient Search**: Search for patients by various criteria (national ID, passport ID, name, etc.)
- **Access Control**: Staff can only access patients from their own hospital
- **External API Integration**: Integrates with external HIS API with automatic fallback to mock data
- **JWT Authentication**: Secure API access using JWT tokens

## Quick Start

### Prerequisites

- Go 1.24+ installed
- PostgreSQL installed and running
- (Optional) Docker and Docker Compose

### 1. Setup Environment Variables

Copy the example environment file and configure it:

```bash
cp example.env .env
```

Edit `.env` with your database credentials:

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=agnos_db
JWT_SECRET=your-secret-key-change-this
HIS_API_BASE_URL=https://hospital-a.api.co.th
```

### 2. Start the Application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 3. Verify It's Running

```bash
curl http://localhost:8080/health
```

## Using Swagger UI

Swagger provides an interactive interface to test all API endpoints.

### Access Swagger

1. Open your browser and navigate to:
   ```
   http://localhost:8080/swagger/index.html
   ```

### How to Use Swagger

#### Step 1: Create a Staff Account

1. Find the **POST /staff/create** endpoint
2. Click "Try it out"
3. The request body is pre-filled with example values
4. Click "Execute" to create a staff account
5. Copy the response (you'll need the username for login)

#### Step 2: Login to Get JWT Token

1. Find the **POST /staff/login** endpoint
2. Click "Try it out"
3. Enter the username and password from Step 1
4. Click "Execute"
5. **Copy the `token` from the response** - you'll need this for patient search

#### Step 3: Authorize for Protected Endpoints

1. Click the **"Authorize"** button at the top right
2. In the dialog, enter: `Bearer YOUR_TOKEN_HERE`
   - Replace `YOUR_TOKEN_HERE` with the token from Step 2
   - **Important**: Include the word "Bearer" followed by a space before your token
3. Click "Authorize" then "Close"

#### Step 4: Search for Patients

1. Find the **GET /patient/search** endpoint
2. Click "Try it out"
3. Enter a patient ID in the `id` field:
   - **Hospital A examples**: `1234567890123`, `9876543210987`, `AB1234567`
   - **Hospital B examples**: `1111222233334`, `4455667788990`
4. Click "Execute"
5. View the patient information in the response

### Available Endpoints

- **POST /staff/create** - Create a new staff account
- **POST /staff/login** - Login and receive JWT token
- **GET /patient/search** - Search for patients (requires JWT authentication)
- **GET /health** - Health check endpoint

## Docker Setup (Optional)

To run the entire stack with Docker:

```bash
# Start all services (PostgreSQL, API, Nginx)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

Access points:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Nginx (reverse proxy): http://localhost:80

## Project Structure

```
AgnosAssigment/
├── cmd/server/main.go       # Application entry point
├── internal/
│   ├── models/             # Data models
│   ├── repositories/       # Database access layer
│   ├── services/           # Business logic
│   ├── controllers/api/    # HTTP handlers
│   ├── middlewares/        # Authentication middleware
│   ├── configs/            # Configuration loader
│   └── utils/              # Utility functions
├── docker-compose.yaml     # Docker services configuration
├── Dockerfile             # Go application container
├── nginx.conf             # Nginx reverse proxy config
└── example.env            # Environment variables template
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

## Notes

- The system automatically falls back to mock data if the external HIS API is unavailable
- Staff can only search for patients from their own hospital
- Patient search supports multiple criteria: national ID, passport ID, name, date of birth, etc.
- All passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
