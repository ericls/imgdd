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

## Features
- **Pluggable storage backends**: S3 (and S3-compatible services like MinIO), local filesystem, WebDAV, and IPFS MFS. Multiple backends can be enabled at once with configurable priority, defined in the database or the config file. Images can be replicated between backends with the `replicate` CLI command.
- **Flexible URL formats**: `canonical` (proxied from the best available backend) and `direct` (routed through a specific backend); per-image delivery chooses the highest-priority enabled backend.
- **Upload pipeline**: automatic EXIF stripping, MIME-type detection with declared/detected mismatch rejection, size limits, and deduplication across stored images. Revisions are linked back to a root image.
- **Optional safe-image check**: delegate NSFW/abuse screening to an external HTTP endpoint.
- **Captcha protection**: Google reCAPTCHA or Cloudflare Turnstile, wired into GraphQL via an `@captchaProtected` directive.
- **Rate limiting** on uploads and other sensitive endpoints.
- **Identity & access control**: email/password auth, organizations, role/permission system with built-in roles, site-owner privilege, and password reset via email.
- **Email backends**: SMTP or a dummy backend for development, with templated messages.
- **GraphQL API** (gqlgen) covering images, viewer, users, organizations, roles/permissions, and storage definitions, with `@isAuthenticated` / `@isSiteOwner` directives.
- **Web client**: React + Tailwind admin & user UI, including a site-admin area for managing users, roles, and storage definitions.
- **Client plugin hooks**: a lightweight `window.registerPlugin` API lets custom JS inject content into named UI slots without rebuilding the frontend.
- **Low resource usage**: a single Go binary serves the API, web client, and image proxying; Instance with 256M memory has comfortably handled 1 TB of traffic and 4.2 M requests per day in production.
- **Operational toggles**: disable new uploads (`ALLOW_UPLOAD`) or new user signups (`ALLOW_NEW_USER`) at runtime.
- **Background cleanup** of orphaned stored images across backends.
- **Deployment-friendly**: single Go binary, official Docker image, TOML + env configuration, optional migrate-on-start, Nix flake for reproducible dev.
- **CLI tooling**: migrations, user/role management, config generation, test email, storage replication, and dev helpers (`gql`, `jet`, `reset-db`, `dev-server` with hot reload).

## Development

### Prerequisites

- Go 1.25+
- Node.js 20 (for the frontend)
- Docker & Docker Compose (for local services)

> If you use nix, you can just do `nix develop`.

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
