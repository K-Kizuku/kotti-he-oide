# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

This is a full-stack hackathon project with three main components:

- **Frontend**: Next.js 15.5.2 React application using TypeScript and Turbopack
- **Server**: Go REST API backend following Layered Architecture + DDD principles
- **Infra**: Infrastructure/deployment configuration

The frontend uses Next.js App Router with TypeScript paths configured (`@/*` maps to `./src/*`).

## Backend Architecture (DDD + Layered)

The Go server follows Domain-Driven Design (DDD) and Clean Architecture principles:

```
server/
├── cmd/server/                    # Application entry point
├── internal/
│   ├── interfaces/                # Interface Layer (外部インターフェース層)
│   │   └── http/
│   │       ├── handler/           # HTTP handlers
│   │       └── dto/               # Data Transfer Objects
│   ├── application/               # Application Layer (アプリケーション層)
│   │   └── usecase/               # Use cases
│   ├── domain/                    # Domain Layer (ドメイン層)
│   │   ├── model/                 # Domain entities
│   │   ├── repository/            # Repository interfaces
│   │   ├── service/               # Domain services
│   │   └── valueobject/           # Value objects
│   └── infrastructure/            # Infrastructure Layer (インフラストラクチャ層)
│       └── persistence/           # Data persistence implementations
└── pkg/
    └── errors/                    # Custom error types
```

### Dependency Flow
- **Interfaces** → **Application** → **Domain** ← **Infrastructure**
- Dependencies point inward (Clean Architecture)
- Infrastructure implements domain interfaces

## Development Commands

### Frontend Development
All frontend commands should be run from the `frontend/` directory:

```bash
cd frontend
pnpm dev          # Start development server with Turbopack
pnpm build        # Build for production with Turbopack  
pnpm start        # Start production server
pnpm lint         # Run ESLint (eslint command only)
```

The frontend runs on http://localhost:3000 by default.

### Server Development
All server commands should be run from the `server/` directory:

```bash
cd server
make run          # Run server from root main.go (port 8080)
make run-cmd      # Run server from cmd/server/main.go
make build        # Build binary to bin/server
make test         # Run tests
make fmt          # Format code
```

The server runs on http://localhost:8080 by default.

### API Endpoints
- `GET /health` - Health check
- `GET /api/users` - Get all users
- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user by ID
- `DELETE /api/users/{id}` - Delete user

## Technology Stack

- **Frontend**: React 19.1.0, Next.js 15.5.2, TypeScript 5+, ESLint 9
- **Backend**: Go 1.25.1, Standard library HTTP server, DDD + Clean Architecture
- **Package Manager**: pnpm (frontend)
- **Build Tool**: Turbopack (Next.js)

## Development Guidelines

- Follow DDD principles: business logic stays in domain layer
- Use value objects for primitive validation (UserID, Email)
- Implement repository pattern for data access abstraction
- Handle errors through custom domain error types
- Use dependency injection for loose coupling