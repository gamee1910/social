# Social

A Golang backend service built with modern API design patterns and best practices.

## Prerequisites

Before getting started, ensure you have the following installed:

- **Go** (1.21 or higher)
- **Git**
- **direnv** - for environment variable management
- **air** - for hot-reloading during development

## Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd social
```

Or download and extract the ZIP folder:

```bash
unzip social.zip
cd social
```

### 2. Install Development Tools

#### Install direnv

Follow the official guide: https://direnv.net/docs/installation.html

**macOS:**
```bash
brew install direnv
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install direnv
```

**Windows (with WSL):**
```bash
sudo apt-get install direnv
```

#### Install air (Hot Reload)

```bash
go install github.com/cosmtrek/air@latest
```

### 3. Configure Environment Variables

**Enable direnv in your shell:**

Open a new terminal in the project directory and run:

```bash
direnv allow
```

This will load environment variables from `.envrc` file.

**Edit `.envrc` for your environment:**

```bash
# Edit the .envrc file
nano .envrc
# or
vim .envrc
```

Add your configuration variables (example):

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=social_db
export SERVER_PORT=8080
export LOG_LEVEL=debug
```

After editing, save and direnv will automatically reload the variables.

### 4. Install Dependencies

```bash
go mod download
go mod tidy
```

## Development

### Running the Application with Hot Reload

Start the development server with automatic restart on code changes:

```bash
air
```

The application will start and watch for file changes in your source code. Any modifications will trigger a rebuild and restart.