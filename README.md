# Social

A Golang backend service built with modern API design patterns and best practices.

## Prerequisites

Before getting started, ensure you have the following installed:

- **Go** (1.21 or higher)
- **Git**
- **Make**
- **direnv** - for environment variable management
- **air** - for hot-reloading during development
- **golang-migrate** - for database migrations

---

## Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd social
```

Or download and extract the ZIP archive:

```bash
unzip social.zip
cd social
```

---

### 2. Install Development Tools

#### Install direnv

Official documentation:

https://direnv.net/docs/installation.html

**macOS**

```bash
brew install direnv
```

**Ubuntu / Debian**

```bash
sudo apt install direnv
```

---

#### Install air (Hot Reload)

```bash
go install github.com/cosmtrek/air@latest
```

Verify:

```bash
air -v
```

---

#### Install golang-migrate

Download the latest binary from:

https://github.com/golang-migrate/migrate/releases

Or install with Go:

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify:

```bash
migrate -version
```

---

### 3. Configure Environment Variables

Enable `direnv` inside the project:

```bash
direnv allow
```

Edit the `.envrc` file to match your local environment.

Example:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=social_db

export SERVER_PORT=8080
export LOG_LEVEL=debug
```

Whenever `.envrc` changes, run:

```bash
direnv allow
```

---

### 4. Install Dependencies

```bash
go mod download
go mod tidy
```

---

## Makefile

This project provides a `Makefile` to simplify common development tasks.

### Show Available Commands

```bash
make help
```

### Run the Application

```bash
make run
```

### Run with Hot Reload

```bash
make watch
```

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

---

## Database Migrations

### Create a Migration

```bash
make migrate-create name=create_users_table
```

This creates:

```
migrations/
├── 000001_create_users_table.up.sql
└── 000001_create_users_table.down.sql
```

---

### Apply All Migrations

```bash
make migrate-up
```

---

### Roll Back All Migrations

```bash
make migrate-down
```

---

### Roll Back One Migration

```bash
make migrate-down-one
```

---

### Check Current Migration Version

```bash
make migrate-version
```

---

### Force a Migration Version

If a migration fails and leaves the database in a dirty state:

```bash
make migrate-force version=3
```

Replace `3` with the migration version you want to force.

---

## Development

Run the development server with hot reload:

```bash
make watch
```

Or without hot reload:

```bash
make run
```

The server will automatically reload whenever source files change when using `make watch`.