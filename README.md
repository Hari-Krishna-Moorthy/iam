# Multi-Tenant IAM Service

A highly scalable, multi-tenant Identity and Access Management (IAM) service built with Go, adhering to strict Domain-Driven Design (DDD) principles.

## Architecture Overview

The project follows a clean architecture with the following layers:

- **Domain Layer (`domain/`):** Contains pure business logic, entities, value objects, and repository interfaces. No external dependencies.
- **Application Layer (`application/`):** Orchestrates use cases and interacts with domain models and repository interfaces.
- **Infrastructure Layer (`infrastructure/`):** Concrete implementations of repositories (GORM for PostgreSQL, Redis for sessions), configurations, and logging.
- **Interfaces Layer (`interfaces/`):** HTTP handlers, middlewares, and routing using `go-chi`.

## Current State: Phase 1 (Domain Definition)

In this phase, the core domain entities and interfaces have been defined.

### Core Domain Entities

- **Tenant:** Represents an organization. Supports multi-tenancy by mapping multiple domains/origins to a single tenant.
- **User:** Represents an identity within a tenant.
- **Role:** A collection of permissions assigned to users.
- **Permission:** A value object following the format `<scope>:<serviceName>:<action>`.
- **Session:** Represents an active user session, managed in Redis for scalability.

### Domain Relationships

```mermaid
classDiagram
    class Tenant {
        +String ID
        +String Name
        +String[] Domains
        +Boolean IsActive
    }

    class User {
        +String ID
        +String TenantID
        +String Username
        +String Email
        +String RoleID
        +Boolean IsActive
    }

    class Role {
        +String ID
        +String TenantID
        +String Name
        +Permission[] Permissions
    }

    class Permission {
        <<Value Object>>
        +String Scope
        +String ServiceName
        +String Action
        +String String()
    }

    class Session {
        +String ID
        +String UserID
        +String TenantID
        +String Role
        +Permission[] Permissions
        +Time ExpiresAt
    }

    Tenant "1" -- "*" User : contains
    Tenant "1" -- "*" Role : defines
    User "*" -- "1" Role : assigned
    Role "1" -- "*" Permission : has
    User "1" -- "*" Session : owns
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
