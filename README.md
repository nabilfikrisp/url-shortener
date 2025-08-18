# URL Shortener

A simple **URL Shortener** service built with **Go (Fiber, GORM, PostgreSQL)**.  
Includes unit tests, integration tests with Docker, and code coverage reporting.

---

## Features

- Shorten long URLs into unique tokens
- Redirect short tokens to original URLs
- PostgreSQL database with GORM
- Unit & integration test coverage
- Docker Compose setup for testing DB
- REST API with 3 endpoints

## Programming Language Choice

I chose **Go** for this project because I will be working with Go in my upcoming role and wanted to take this opportunity to learn the language and its ecosystem.

## Design Decisions & Assumptions

- **Repository Pattern**: I implemented the repository pattern which commonly used in Go applications to separate data access logic from business logic. As I'm still learning Go idioms, some patterns may reflect influences from my previous TypeScript experience.
- **Database Choice**: Selected PostgreSQL with GORM as I have the most experience with PostgreSQL from previous projects.
- **My Thoughts About This Project**: While this serves as both a technical assessment and a learning exercise, I've aimed for best practices and functional requirements. Both the application and test cases could have been implemented much faster if i use TypeScript with Express/Hono/NestJS, but the majority of development time was invested in learning Go syntax, patterns, and best practices.

---

# URL Shortener – Go API

## Getting Started (Local Development)

### Prerequisites

- Go **1.24.3** (or newer)
  ```bash
  go version
  # go version go1.24.3 windows/amd64
  ```

* PostgreSQL (either installed locally, or run via Docker)

---

### Clone the Repository

```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener/go-api
```

---

### Setup Environment

Copy the example env file and adjust values:

```bash
cp .env.example .env
```

or simply manually create .env

---

### Start PostgreSQL

- **Option A: Use Local PostgreSQL** (make sure it matches the config in `.env`)
- **Option B: Use Docker (recommended if no local Postgres installed):**

  ```bash
  docker compose -f docker-compose.dev.yml up -d dev-db
  ```

  or

  ```bash
  make dev-db-up
  ```

  > This runs PostgreSQL on **port 5433**.

  > Please adjust the .env file if using this option

- **If Using Option B: Shut Down DB by doing**

  ```bash
  docker compose -f docker-compose.dev.yml down -v
  ```

  or

  ```bash
  make dev-db-down
  ```

---

### Install Dependencies

```bash
go mod tidy
```

---

### Run the Application

```bash
go run cmd/server/main.go
```

The server will start on (port based on env):
👉 [http://localhost:3001](http://localhost:3001)

---

## Getting Started Another Option (DOCKER)

### Prerequisites

- Docker
  ```bash
  docker --version
  # Docker version 28.3.2, build 578ccf6
  ```

---

### Clone the Repository

```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener/go-api
```

---

### Start Right Away

```bash
docker compose up --build
```

---

## Endpoints

### Create Short Token

**POST** `/shorten`

**Request body:**

```json
{
  "url": "http://example.com"
}
```

---

### Redirect Short Token

**GET** `/:shortToken`

**Example:**

```
GET /abc123
```

→ Redirects to the original URL.

---

### Get Original URL (JSON)

**GET** `/api/urls/:shortToken`

**Example:**

```
GET /api/urls/abc123
```

**Response:**

```json
{
  "message": "URL retrieved successfully",
  "data": {
    "id": 1,
    "short_token": "8701e897618bc220",
    "original": "https://claude.ai/",
    "click_count": 2,
    "created_at": "2025-08-17T16:28:04.763986Z",
    "updated_at": "2025-08-17T16:28:04.763986Z"
  }
}
```

---

## Testing

This project separates **unit tests** and **integration tests**, although both can be run together.

- Unit tests run without any external dependencies
- Integration tests require a running Postgres database.

---

### 🗄️ Setup Test Database

- If you already have a database set up manually, just update the values in **`.env.test`**.
- If not, you can spin up a Postgres instance using Docker.

#### Start Docker test database

```bash
docker compose -f docker-compose.test.yml up -d test-db
# or
make test-db-up
```

#### Stop & remove Docker test database

```bash
docker compose -f docker-compose.test.yml down -v
# or
make test-db-down
```

> When using the Docker Postgres test DB, simply copy the values from **`.env.test.example`** into **`.env.test`**.

---

### ⚙️ Setup Environment

Create your **`.env.test`** file by either:

```bash
cp .env.test.example .env.test
```

or manually writing one yourself.

---

### Run All Tests

```bash
make test
# or
go test ./... -v
```

### Run Unit Tests

```bash
make test-unit
# or
go test ./tests/unit/... -v
```

### Run Integration Tests

```bash
make test-integration
# or
go test ./tests/integration/... -v
```

---

## Code Coverage

Generate a coverage report:

```bash
make coverage
```

OR

Generate coverage

```bash
go test ./... -coverprofile=coverage.out
```

then

```bash
go tool cover -html="coverage.out"
```

---
