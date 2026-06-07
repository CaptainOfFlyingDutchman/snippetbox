# Snippetbox

A web application written in Go.

## Prerequisites

Before running the application, set up the database by following the guide below:

- [Database Setup](docs/database.md)

## Running the Application

Start the web server:

```bash
cd snippetbox
go run ./cmd/web
```

The application will be available at:

```text
http://localhost:4000
```

## Development Tools

The database setup guide also includes configuration instructions for:

- DataGrip
- WebStorm
- Other MySQL-compatible clients

## Generating TLS certificate

```bash
mkdir tls
go run <GO_Install_Path>/<GO_Version>/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```