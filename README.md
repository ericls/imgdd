# IMGDD

[![Backend test](https://github.com/ericls/imgdd/actions/workflows/backend.yaml/badge.svg?branch=main)](https://github.com/ericls/imgdd/actions/workflows/backend.yaml)
[![GitHub Release](https://img.shields.io/github/v/release/ericls/imgdd)](https://github.com/ericls/imgdd/releases/)
[![Docker Image Version](https://img.shields.io/docker/v/ericls/imgdd)](https://hub.docker.com/r/ericls/imgdd)

## Introduction
IMGDD is a simple, self-hostable image hosting program.

It powers an image hosting project that began in early 2023. After some users requested the source code, I decided to make it available following a bit of refactoring and cleanup.

As of February 2025, the project handles over 1 TB of traffic and 4.2 million requests per day. Since continued growth may force me to stop accepting new images, I open-sourced IMGDD so that anyone can run their own instance and serve their images.

## Live Instances
- [imgdd.com](https://imgdd.com)

## Development

### Prerequisites

- Go 1.25+
- Node.js 20 (for the frontend)
- Docker & Docker Compose (for local services)

### Setup

1. Copy the environment file and adjust if needed:
   ```bash
   cp .env.template .env
   ```

2. Start local services (PostgreSQL, Redis, MinIO ancient version):
   ```bash
   docker compose up -d
   ```

3. Run database migrations:
   ```bash
   go run . migrate
   ```

4. Populate built-in roles:
   ```bash
   go run . populate-built-in-roles
   ```

5. Create a user:
   ```bash
   go run . create-user --email you@imgdd.com --password yourpassword --is-site-owner
   ```

### Building & Running

```bash
# Build the frontend
cd web_client && pnpm install && pnpm build && cd ..

# Start the server
go run . serve

# Start with auto-reload (requires air)
go run . dev-server
```

### Code Generation

Dev commands are available when running from source:

```bash
# Regenerate GraphQL code (gqlgen)
go run . gql

# Regenerate Jet ORM code (requires running PostgreSQL)
go run . jet

# Regenerate frontend GraphQL types
cd web_client && pnpm gen
```

### Database Commands

```bash
# Run migrations
go run . migrate

# Migrate to a specific version
go run . migrate --version 3

# Create a new migration file
go run . make-migration --name add_some_table

# Reset the database (dev only)
go run . reset-db
```

All commands accept a `-c` / `--config` flag to specify the config file path, and a `--log-level` flag (default: `info`).

## Deployment
See [Start guide in the Wiki](https://github.com/ericls/imgdd/wiki/Start-guide).

## FAQ
See [FAQ in the Wiki](https://github.com/ericls/imgdd/wiki/FAQ).
