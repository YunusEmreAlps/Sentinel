<br />
<p>
    <a href="" target="_blank">
      <img
        src="https://github.com/YunusEmreAlps/Sentinel/blob/master/assets/sentinel.png?raw=true"
        alt="Sentinel"
        width="100%"
      />
    </a>
</p>

In the realm of digital security, "Sentinel" stands tall as an open-source powerhouse, meticulously crafted to monitor and control certificate expiration dates with unparalleled precision. With the release of version 1.0.0, Sentinel has ascended to new heights, now adept at retrieving critical certificate information from list and promptly alerting designated teams through comprehensive notifications, ensuring your systems remain secure.

## Sentinel Meaning for this project

Vigilant Protector: Embodying the essence of a sentinel—a vigilant guardian standing watch—Sentinel epitomizes the system's core mission of safeguarding your digital infrastructure by actively monitoring and controlling certificate expiration dates.

Stalwart and Resilient: The name Sentinel underscores the system's strength and resilience. Like an unwavering guard, it stands firm against the risks associated with expired certificates, ensuring the robustness and reliability of your security infrastructure.

Futuristic Connotations: "Sentinel" is a term associated with advanced technology and cutting-edge security measures. By adopting this name, the project signals a commitment to staying ahead in the dynamic landscape of digital security, aligning with futuristic and technological connotations.

With Sentinel 1.0.0, rest assured that your certificate expiration dates are not merely managed but guarded with utmost precision, reflecting a dedication to excellence in certificate control services.

## Prerequisites

- Go 1.15+
- PostgreSQL
- Docker & Docker Compose
- SMTP

## Quick start

We can run this **Sentinel** project with or without Docker. Here, I am providing both ways to run this project.

- Clone this project

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone ...

# Move to the project root directory
cd Sentinel
```

### Run without Docker

Run the following command to execute the Go program:

- Make sure PostgreSQL is running and accessible with the credentials you provided in the .env file.
- Open a terminal or command prompt and navigate to the root directory of your project.
- Create a file `.env` similar to `.env.example` at the **/config directory** with your configuration.
- Install `go` if not installed on your machine.
- Install `PostgreSQL` if not installed on your machine.
- Important: Open the `.env` file and modify the values of `DB_HOST`, `DB_USER`, and `DB_PASSWORD` to match your PostgreSQL configuration. Update any other configuration variables if necessary.
- Run `go run main.go`.

### Run with Docker

- Create a file `.env` similar to `.env.example` at the **/config directory** with your configuration.
- Install Docker and Docker Compose.
- Run `docker-compose up -d`.

## Project structure

### `config`

Configuration. First, `config.yml` is read, then environment variables overwrite the yaml config if they match.
The config structure is in the `config.go`.
The `env-required: true` tag obliges you to specify a value (either in yaml, or in environment variables).

Reading the config from yaml contradicts the ideology of 12 factors, but in practice, it is more convenient than
reading the entire config from ENV.
It is assumed that default values are in yaml, and security-sensitive variables are defined in ENV.

### `main.go`

Core of this project (DB Connection, Excelize, Filtering and more...).

## Major Packages used in this project

- **postgeSQL go driver**: The Official Golang driver for PostgreSQL.
- **gorm**: The fantastic ORM library for Golang, aims to be developer friendly.
- **viper**: For loading configuration from the `.env` file. Go configuration with fangs. Find, load, and unmarshal a configuration file in JSON, TOML, YAML, HCL, INI, envfile, or Java properties formats.
- **bcrypt**: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
- **testify**: A toolkit with common assertions and mocks that plays nicely with the standard library.
- Check more packages in `go.mod`.

## Contributing

We welcome contributions to Orion! To contribute to the project, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Make your changes and ensure that the tests pass.
- Commit your changes and push them to your fork.
- Submit a pull request to the main repository, describing your changes in detail.
- Please review the Contribution Guidelines for more information.
