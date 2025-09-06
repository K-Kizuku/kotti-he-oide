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

### Infrastructure Development
All infrastructure commands should be run from the `infra/` directory:

```bash
cd infra
terraform init      # Initialize Terraform
terraform plan       # Preview infrastructure changes
terraform apply      # Apply infrastructure changes
terraform destroy    # Destroy infrastructure resources
```

### API Endpoints
- `GET /health` - Health check
- `GET /api/users` - Get all users
- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user by ID
- `DELETE /api/users/{id}` - Delete user

## Infrastructure Architecture

The infrastructure is deployed on AWS using Terraform and follows a containerized approach with ECS Fargate:

```
infra/
├── alb.tf                  # Application Load Balancer configuration
├── ecr.tf                  # Elastic Container Registry
├── ecs_cluster.tf          # ECS Cluster setup
├── ecs_services_api.tf     # API service configuration
├── ecs_services_web.tf     # Web service configuration
├── outputs.tf              # Terraform outputs
├── providers.tf            # AWS and Random providers
├── rds.tf                  # RDS PostgreSQL database
├── s3.tf                   # S3 bucket configuration
├── security.tf             # Security groups and IAM roles
├── variables.tf            # Input variables
├── versions.tf             # Terraform version constraints
└── vpc.tf                  # VPC and networking
```

### Infrastructure Components

- **Compute**: ECS Fargate cluster hosting containerized API and Web services
- **Load Balancing**: Application Load Balancer (ALB) with path-based routing
- **Container Registry**: ECR repositories for API and Web images
- **Database**: RDS PostgreSQL instance
- **Storage**: S3 bucket for static assets
- **Networking**: Custom VPC with public/private subnets
- **Security**: Security groups and IAM roles for least privilege access

### Deployment Process

1. Build and push container images to ECR
2. Configure `terraform.tfvars` with required variables
3. Deploy infrastructure using Terraform
4. Access services via ALB DNS name:
   - `/` - Web service (Next.js frontend)
   - `/api/*` - API service (Go backend)

## Technology Stack

- **Frontend**: React 19.1.0, Next.js 15.5.2, TypeScript 5+, ESLint 9
- **Backend**: Go 1.25.1, Standard library HTTP server, DDD + Clean Architecture
- **Infrastructure**: AWS ECS Fargate, ALB, RDS PostgreSQL, ECR, S3, Terraform ~> 5.0
- **Package Manager**: pnpm (frontend)
- **Build Tool**: Turbopack (Next.js)

## Development Guidelines

- Follow DDD principles: business logic stays in domain layer
- Use value objects for primitive validation (UserID, Email)
- Implement repository pattern for data access abstraction
- Handle errors through custom domain error types
- Use dependency injection for loose coupling