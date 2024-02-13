# Go-Micro

<p>
  Project build to play around with microservices.
</p>

## Development requirements
- [Go](https://go.dev/dl/) installation, the application was developed using Go v1.21.6

- [GNU Make](https://www.gnu.org/software/make/) installation, the application was developed using GNU Make v4.3

- [Docker](https://www.docker.com/products/docker-desktop/) installation, the application was developed using Docker v24.0.7


## Running the app

### Starting up
1. `make up_build` Build the service binaries and start them up
2. `make start` Build front-end binary and start it up

### Shutting down
1. `make stop` Stop the front-end and removes the binary
2. `make down` Stops services and removes their binaries

## Database connection
- **MongoDB:** `mongodb://admin:password@localhost:27017/logs?authSource=admin&tls=false`

- **PostgreSQL:** `host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable`
