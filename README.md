# 🚀 Golang REST API

This is an **example project** built with **Go** and **PostgreSQL**, supporting **Swagger**, **Docker**, and **Migrations**.  
It serves as a **template** for building RESTful APIs with Golang.  

---

## 📦 Prerequisites

Before diving in, make sure you have the following installed on your system:

- **[Git](https://git-scm.com/downloads)**  
  Version control system for tracking changes.

- **[Go v1.22+](https://go.dev/dl)**  
  The latest version of the Go programming language.

- **[PostgreSQL](https://www.postgresql.org/download/)**  
  A powerful, open-source object-relational database system.  
  **OR** use **[Docker](https://hub.docker.com/_/postgres)** to run PostgreSQL in a container.

- **[Golang Migrate](https://github.com/golang-migrate/migrate/tree/master)**  
  A migration tool for managing database schema changes.

- **[Swagger](https://github.com/swaggo/swag/)**  
  For generating interactive API documentation.  
  Install it with:
  ```sh
  go install github.com/swaggo/swag/cmd/swag@v1.16.4
  ```

- **Docker & Docker Compose**  
  For containerization and managing multi-container Docker applications.

---

## 🚀 Quick Start

### 1. Generate Swagger Documentation

Swagger documentation is crucial for understanding and testing your API endpoints. To generate the documentation, run:

```sh
swag init -g cmd/api/main.go
```

This command scans your code for Swagger annotations and creates the necessary documentation files.

### 2. Start the API Server

To launch your API server in development mode, execute:

```sh
go run cmd/api/main.go
```

Your API is now running, and you can start making requests.

---

## 🐳 PostgreSQL with Docker

If you prefer running PostgreSQL in a container, follow these steps:

1. **Copy the Environment File**

   Create a new environment configuration file by copying the example:

   ```sh
   cp docker/env.example docker/.env
   ```

2. **Customize Your Environment Variables**

   Open the `docker/.env` file and update the variables (e.g., database name, user, password) to match your desired configuration.

3. **Start the PostgreSQL Container**

   Launch the container using Docker Compose:

   ```sh
   docker compose -f docker/docker-compose.yml -p golang-rest-api-infra up -d
   ```

4. **Verify the Container Status**

   Check that your container is running properly:

   ```sh
   docker ps -a
   ```

---

## 📂 Database Migrations - Keep It Fresh!

Managing your database schema is critical. Use Golang Migrate to handle migrations.

### Create a New Migration

Generate a new migration file by running:

```sh
make new_migration name={migration_name}
```

Replace `{migration_name}` with a descriptive name for your migration.

### Apply Migrations

Before running migrations, set your database URL:

```sh
export DB_URL=postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}?sslmode=disable
```

Then apply all pending migrations:

```sh
make migrateup
```

This will update your database schema to match the latest version.

---

## 🔹 Seed Initial Data

Seeding data is essential for testing and initial setup.

1. **Copy the Seeding Environment File**

   ```sh
   cp script/seed_user/env.example script/seed_user/.env
   ```

2. **Update the `.env` File**

   Edit the file with your specific database configuration details.

3. **Prepare the CSV File**

   Fill in the CSV file with the initial user data.

4. **Run the Seeding Script**

   Populate your database with the seed data:

   ```sh
   go run script/seed_user/seed_user.go
   ```

---

## 🔐 Generate RSA Keys for JWT - Lock It Down!

Secure your API with JWT-based authentication by generating RSA keys.

### Generate the Private Key

Create a secure private key:

```sh
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```

### Generate the Public Key

Extract the public key from the private key:

```sh
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

These keys will be used for signing and verifying JWT tokens.

