# Command Line Tools

This directory contains various command-line tools and entry points for the Katseye API application.

## Directory Structure

```
cmd/
├── api/                  # Main API application entry point
│   └── main.go           # Initializes and runs the HTTP server
├── migrations/           # Database migration scripts
│   └── migrate_product_partner/ # Migration for product partner data
└── seed_user/            # User seeding utility
    └── main.go           # Creates initial user accounts
```

## Tools

### API Server (`api/`)

The main application entry point that initializes and runs the HTTP server. It handles graceful shutdown on system signals.

**Usage:**
```
go run cmd/api/main.go
```

### Seed User (`seed_user/`)

A utility to create user accounts in the database. This tool is useful for creating initial admin accounts or service accounts.

**Usage:**
```
go run cmd/seed_user/main.go -email=user@example.com -password=securepassword -role=admin
```

**Parameters:**
- `-email`: Email address for the user (required)
- `-password`: Password in plain text (will be hashed) (required)
- `-active`: Whether the user should be active (default: true)
- `-role`: User role (admin, manager, user) (default: user)
- `-profile_type`: Type of profile (service_account, partner_manager, consumer) (default: service_account)

### Migrations (`migrations/`)

Contains database migration scripts for schema changes and data transformations.

#### Product Partner Migration (`migrate_product_partner/`)

Migration script for product partner data.