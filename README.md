# Multi-Tenant IAM Service

A highly scalable, multi-tenant Identity and Access Management (IAM) service built with Go, adhering to strict Domain-Driven Design (DDD) principles.

## Local Development

To run the service locally with a real database and Redis:

1.  **Start Services:** Use Docker Compose to spin up PostgreSQL and Redis.
    ```bash
    docker-compose up -d
    ```

2.  **Configuration:** The application uses a `.env` file for configuration. A default `.env` has been provided with the following settings:
    ```text
    DB_URL=postgres://postgres:postgres@localhost:5432/iam?sslmode=disable
    REDIS_URL=localhost:6379
    PORT=8080
    ```

3.  **Run the Server:**
    ```bash
    go run cmd/server/main.go
    ```

4.  **Run Tests:**
    ```bash
    go test ./...
    ```

## Architecture Overview

The project follows a clean architecture with the following layers:

- **Domain Layer (`domain/`):** Contains pure business logic, entities, value objects, and repository interfaces. No external dependencies.
- **Application Layer (`application/`):** Orchestrates use cases and interacts with domain models and repository interfaces.
- **Infrastructure Layer (`infrastructure/`):** Concrete implementations of repositories (GORM for PostgreSQL, Redis for sessions), configurations, and logging.
- **Interfaces Layer (`interfaces/`):** HTTP handlers, middlewares, and routing using `go-chi`.

## Current State: Phase 2 (Infrastructure Implementation)

In this phase, the infrastructure layer has been implemented, providing concrete data access logic.

### Infrastructure Components

- **GORM Models:** PostgreSQL-specific models with struct tags for `Tenant`, `User`, and `Role`.
- **PostgreSQL Repositories:** Concrete implementations of `domain` interfaces using GORM.
- **Redis Session Repository:** Manages active user sessions in Redis, supporting multiple sessions per user and strict session payloads.
- **Mapping:** Strict separation between pure `domain` entities and `infrastructure` models using mapper functions (`ToDomain` / `FromDomain`).

### Database Schema (Conceptual)

```mermaid
erDiagram
    TENANT ||--o{ USER : contains
    TENANT ||--o{ ROLE : defines
    USER }|--|| ROLE : assigned
    SESSION }|--|| USER : belongs_to

    TENANT {
        uuid id PK
        string name
        text[] domains
        boolean is_active
    }

    USER {
        uuid id PK
        uuid tenant_id FK
        string username
        string email
        string password_hash
        uuid role_id FK
    }

    ROLE {
        uuid id PK
        uuid tenant_id FK
        string name
        text[] permissions
    }

    SESSION {
        string id PK
        string user_id
        string tenant_id
        string role
        text[] permissions
    }
```

## Directory Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── domain/
│   ├── permission/
│   │   └── value_objects.go
│   ├── role/
│   │   ├── entity.go
│   │   └── repository.go (interface in entity.go)
│   ├── session/
│   │   ├── entity.go
│   │   └── strategy.go (interface in entity.go)
│   ├── tenant/
│   │   ├── entity.go
│   │   └── repository.go (interface in entity.go)
│   └── user/
│       ├── entity.go
│       └── repository.go (interface in entity.go)
├── application/
│   ├── auth/
│   ├── session/
│   ├── tenant/
│   └── user/
├── infrastructure/
│   ├── config/
│   ├── logger/
│   └── persistence/
│       ├── gorm/
│       │   ├── models/
│       │   └── repositories/
│       └── redis/
│           └── repositories/
├── interfaces/
│   └── http/
│       ├── handlers/
│       ├── middleware/
│       └── router.go
└── README.md
```
