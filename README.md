# Multi-Tenant IAM Service

A highly scalable, multi-tenant Identity and Access Management (IAM) service built with Go, adhering to strict Domain-Driven Design (DDD) and Test-Driven Development (TDD) principles.

## Features

- **Origin-Based Multi-Tenancy:** Seamlessly identifies tenants based on request `Origin` or `Host` headers.
- **Robust Authentication (JWT & Strategies):** Uses a Strategy Pattern for extensible auth methods. Pre-configured with secure Bcrypt password hashing and issues signed JWTs.
- **Redis Session Management:** Tracks multiple active sessions per user, allowing for instant invalidation and centralized state management.
- **Downstream Request Hydration:** Middleware validates JWTs and injects `X-User-ID`, `X-Tenant-ID`, `X-Role`, and `X-Permissions` into downstream HTTP headers.
- **Advanced Rate Limiting:** Multi-tenant rate limiting backed by Redis. Supports 4 algorithms (Fixed Window, Sliding Window, Leaky Bucket, Token Bucket) configured per-tenant.
- **Asynchronous Audit Logging:** Captures all authentication and management actions into PostgreSQL via non-blocking goroutines.
- **RBAC & Group Management:** Dynamic APIs for managing Roles, Permissions, and hierarchical User Groups (Many-to-Many relationships).
- **Password Policy Enforcement:** Tenant-specific configurable password policies (length, numbers, special characters, uppercase).
- **Graceful Shutdown:** Safely drains active HTTP requests and background tasks upon receiving OS termination signals.

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

- **Domain Layer (`domain/`):** Contains pure business logic, entities, value objects, and repository interfaces. No external dependencies or GORM tags.
- **Application Layer (`application/`):** Orchestrates use cases (Auth, Users, Roles, Rate Limits) and interacts with domain models.
- **Infrastructure Layer (`infrastructure/`):** Concrete implementations of repositories (GORM for PostgreSQL, Redis), Configurations, Migrations, and JWT Providers.
- **Interfaces Layer (`interfaces/`):** HTTP handlers, middlewares (Tenant, Auth, Audit, RateLimit), and routing using `go-chi`.

### Database Schema (Conceptual)

```mermaid
erDiagram
    TENANT ||--o{ USER : contains
    TENANT ||--o{ ROLE : defines
    TENANT ||--o{ GROUP : organizes
    TENANT ||--o| RATE_LIMIT_CONFIG : configures
    TENANT ||--o| PASSWORD_POLICY : enforces
    TENANT ||--o{ AUDIT_LOG : tracks
    USER }|--|| ROLE : assigned
    USER }|--o{ GROUP_USERS : belongs_to
    GROUP ||--o{ GROUP_USERS : includes
    GROUP ||--o{ GROUP_ROLES : has
    ROLE ||--o{ GROUP_ROLES : assigned_to

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

    GROUP {
        uuid id PK
        uuid tenant_id FK
        string name
    }

    RATE_LIMIT_CONFIG {
        uuid tenant_id PK
        string algorithm
        int limit
        int burst
    }
    
    PASSWORD_POLICY {
        uuid tenant_id PK
        int min_length
        boolean require_number
    }

    AUDIT_LOG {
        uuid id PK
        uuid tenant_id FK
        string action
        string resource
    }
```

## Directory Structure Highlights

```text
.
├── cmd/server/main.go                   # Bootstrapping & Graceful Shutdown
├── domain/                              # Pure Entities & Interfaces
│   ├── ratelimit/
│   ├── audit/
│   ├── user/
│   ├── role/
│   ├── tenant/
│   └── session/
├── application/                         # Use Cases
│   ├── auth/strategies/
│   ├── ratelimit/
│   ├── role/
│   ├── user/
│   └── tenant/
├── infrastructure/
│   ├── auth/                            # JWT Provider
│   ├── config/                          # Dotenv Loader
│   └── persistence/
│       ├── registry.go                  # Repo Container
│       ├── gorm/migrations/             # Auto-Migrations
│       ├── gorm/models/                 # DB Structs with Tags
│       ├── gorm/repositories/           # PostgreSQL Repos
│       └── redis/repositories/          # Redis Repos
└── interfaces/http/
    ├── handlers/                        # HTTP Controllers
    ├── middleware/                      # Interceptors (Audit, Auth, Tenant, RateLimit)
    └── router.go                        # Chi Router Setup
```
