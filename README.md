<!-- SENTINEL -->

<br />

<p align="center">
  <a href="assets/sentinel.png" target="_blank">
    <picture>
      <img src="assets/sentinel.png" alt="Sentinel Logo" align="center">
    </picture>
  </a>
</p>

# Sentinel (Certificate Expiration Checker)

In the realm of digital security, "Sentinel" stands tall as an open-source powerhouse, meticulously crafted to monitor and control certificate expiration dates with unparalleled precision. With the release of version 1.0.0, Sentinel has ascended to new heights, now adept at retrieving critical certificate information from list and promptly alerting designated teams through comprehensive notifications, ensuring your systems remain secure.

> This project support Openshift Serverless Architecture (Tested and works perfectly)

## Table of Contents

- [Sentinel (Certificate Expiration Checker)](#sentinel-certificate-expiration-checker)
  - [Table of Contents](#table-of-contents)
  - [Meaning of Sentinel](#meaning-of-sentinel)
  - [Development Guide](#development-guide)
    - [1. Prerequisites](#1-prerequisites)
    - [2. Installation](#2-installation)
    - [4. Running the Application](#4-running-the-application)
      - [Run without Docker](#run-without-docker)
      - [Run with Docker](#run-with-docker)
    - [3. Environment Variables](#3-environment-variables)
  - [Project Structure](#project-structure)
    - [Key Components](#key-components)
  - [Major Packages](#major-packages)
  - [Contributing](#contributing)
  - [API Documentation](#api-documentation)
    - [Accessing Swagger Documentation](#accessing-swagger-documentation)
    - [API Endpoints](#api-endpoints)
    - [Example API Response](#example-api-response)
  - [License](#license)


## Meaning of Sentinel

**Mythological & Historical Context:**
In ancient mythology and military history, a sentinel represents a watchful guardian or sentry posted to keep watch and warn of approaching danger. These vigilant protectors were entrusted with the critical responsibility of maintaining constant surveillance, ensuring the safety of their domain through unwavering attention and immediate alert systems.

**For This Project:**
Sentinel embodies the essence of a vigilant guardian standing watch, epitomizing the system's core mission of safeguarding your digital infrastructure by actively monitoring and controlling certificate expiration dates. Like its mythological counterpart, this digital sentinel maintains constant vigilance over your certificates, providing early warnings and comprehensive protection against security vulnerabilities.

## Development Guide

Sentinel is designed to be a robust and efficient tool for managing certificate expiration dates. Below is a guide to help you set up and run the project, whether you prefer using Docker or running it directly on your machine.

### 1. Prerequisites

- [Go 1.21](https://go.dev/dl/) (The project is developed using Go 1.21.1)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Postman](https://www.postman.com/downloads/)
- [Git](https://git-scm.com/downloads)

> **Note**: The project uses Go 1.21, so make sure to have the correct version installed. You can check the Go version using the following command:

```bash
go version
```

### 2. Installation

We can run this **Sentinel** project with or without Docker. Here, I am providing both ways to run this project.

- Clone the repository

```bash
git clone <repository-url>
```

- Go to your workspace

```bash
cd <path-to-your-workspace>
```

- Navigate to the project directory

```bash
cd Sentinel
```

- Create a file `.env` similar to `.env.example` at the **/config directory** with your configuration.

> **Note**: In this version, No need to database connection. The project uses a REST API to fetch certificate information from a list, so you don't need to set up a PostgreSQL database.

### 4. Running the Application

Choose between two deployment options based on your requirements:

#### Run without Docker

Follow these steps to run Sentinel directly on your machine:

1. **Ensure Prerequisites**: Make sure Go 1.21+ is installed and accessible
2. **Environment Setup**: Create a `.env` file in the `/config` directory based on `sample.env.yaml`
3. **Install Dependencies**: Run `go mod tidy` to download all required packages
4. **Start Application**: Execute `go run main.go` from the project root

> **Note**: For Version 1 (CRON-based), the application will run as a background service. For Version 2 (REST API), it will start a web server.

#### Run with Docker

For containerized deployment:

1. **Environment Setup**: Create a `.env` file in the `/config` directory
2. **Docker Deployment**: Run `docker-compose up -d`
3. **Verify**: Check container status with `docker-compose ps`

> **Note**: Docker Compose will handle all dependencies and networking automatically.

### 3. Environment Variables

- Create a file `.env` similar to `.env.example` at the **/config directory** with your configuration.
- The `.env` file contains environment variables that are used to configure the application. It is essential to set these variables correctly for the application to function properly.

## Project Structure

```text
Sentinel/
├── assets/                 # Project assets and images
│   └── sentinel.png
├── config/                 # Configuration files
│   ├── config.go          # Configuration structure
│   └── .env.example       # Environment variables template
├── internal/               # Internal application packages
│   ├── handlers/          # HTTP request handlers
│   │   ├── get_all_expirations.go
│   │   ├── get_certificate_info.go
│   │   ├── handlers.go
│   │   └── list_certificates.go
│   └── models/            # Data models
│       ├── log.go
│       └── mail.go
├── pkg/                   # Public packages
│   ├── db/               # Database connections
│   │   ├── postgres/
│   │   └── redis/
│   ├── logger/           # Logging utilities
│   ├── mail/             # Email functionality
│   ├── metric/           # Metrics and monitoring
│   ├── parseHtml/        # HTML template parsing
│   ├── templates/        # HTML templates
│   └── utils/            # Utility functions
├── Dockerfile            # Docker container definition
├── docker-compose.yml    # Docker Compose configuration
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point
└── README.md            # Project documentation
```

### Key Components

- **`config/`**: Configuration management with YAML and environment variable support
- **`internal/handlers/`**: HTTP API endpoints for certificate operations
- **`internal/models/`**: Data structures for logs and email notifications
- **`pkg/`**: Reusable packages for database, logging, metrics, and utilities
- **`main.go`**: Application bootstrap with server initialization

## Major Packages

- **gin-gonic/gin**: High-performance HTTP web framework written in Go. Features include JSON validation, error management, middleware support, and excellent performance.
- **gorm**: The fantastic ORM library for Golang with support for PostgreSQL, MySQL, SQLite, and SQL Server. Aims to be developer friendly with features like auto migration, associations, hooks, and transactions.
- **gorm.io/driver/postgres**: Official PostgreSQL driver for GORM ORM.
- **viper**: Go configuration with fangs. Find, load, and unmarshal configuration files in JSON, TOML, YAML, HCL, INI, envfile, or Java properties formats.
- **go-redis/redis**: Type-safe Redis client for Golang with support for Redis 6 commands, pipelining, transactions, and Lua scripting.
- **excelize**: Go library for reading and writing Microsoft Excel™ (XLSX) files. Used for generating Excel reports of certificate information.
- **logrus**: Structured logger for Go, completely API compatible with the standard library logger.
- **gin-contrib/cache**: Gin middleware for HTTP caching with multiple storage backends including Redis and in-memory cache.
- **swaggo**: Automatically generate RESTful API documentation with Swagger 2.0 for Go. Includes gin-swagger for Gin framework integration.
- **gomail**: Simple and efficient package for sending emails. Supports HTML emails, attachments, and embedded files.
- **jaeger**: Distributed tracing system for monitoring and troubleshooting microservices-based distributed systems.
- **prometheus**: Monitoring system and time series database with powerful query language and alerting capabilities.
- **opentracing**: Vendor-neutral APIs and instrumentation for distributed tracing.
- Check more packages in `go.mod`.

## Contributing

We welcome contributions to Sentinel! To contribute to the project, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Make your changes and ensure that the tests pass.
- Commit your changes and push them to your fork.
- Submit a pull request to the main repository, describing your changes in detail.
- Please review the Contribution Guidelines for more information.

## API Documentation

Sentinel provides comprehensive API documentation through Swagger/OpenAPI specification for easy integration and testing.

### Accessing Swagger Documentation

1. **Start the Application**: Ensure Sentinel is running locally or on a server
2. **Open Browser**: Navigate to `http://localhost:8080/sentinel/swagger/index.html`
3. **Explore**: Use the interactive interface to test API endpoints

### API Endpoints

- **GET `/certificates`**: List all monitored certificates
- **GET `/certificates/{domain}`**: Get certificate information for a specific domain
- **GET `/expirations`**: Get certificates approaching expiration
- **POST `/check`**: Manually trigger certificate check

### Example API Response

```json
{
  "status": true,
  "intent": "cld:::sentinel:::/certificates/:domain",
  "message": [
    {
      "version": 3,
      "serial_number": "6523920412297345150429844430373140538",
      "subject": "CN=*.netlify.app,O=Netlify\\, Inc,L=San Francisco,ST=California,C=US",
      "issuer_subject": "CN=DigiCert Global G2 TLS RSA SHA256 2020 CA1,O=DigiCert Inc,C=US",
      "domain": "yunusemrealpu.netlify.app",
      "port": 443,
      "common_name": "*.netlify.app",
      "organization": "Netlify, Inc,",
      "issued_on": "2025-01-31T00:00:00Z",
      "expires_on": "2026-03-03T23:59:59Z",
      "certificate_data": "-----BEGIN CERTIFICATE----\nMIIGFzCCBP\n...",
      "signature_algorithm": "SHA256-RSA",
      "subject_key_id": "3e6abe6e25ac1210abbef1eba7a9bc6d887d548f",
      "authority_key_id": "748580c066c7df37decfbd2937aa031dbeedcd17",
      "is_ca": false,
      "issuer": "DigiCert Global G2 TLS RSA SHA256 2020 CA1",
      "is_expired": false,
      "message": "Certificate will expire in 238 days.",
      "status": 1
    }
  ]
}
```

<!-- Image List -->

![Swagger 1](assets/1.png)
![Swagger 2](assets/2.png)
![Swagger 3](assets/3.png)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
