# Nuitee

A Go-based hotel data management and review aggregation system with REST API capabilities.

## Overview

Nuitee is a comprehensive hotel management system that provides:
- Hotel data storage and management
- Review aggregation and analysis
- RESTful API for data access
- Data synchronization from external sources
- Search and filtering capabilities

## Features

- **Hotel Management**: Store and manage hotel information including ratings, amenities, and contact details
- **Review System**: Aggregate and store hotel reviews from multiple sources
- **REST API**: Full CRUD operations via HTTP endpoints
- **Data Sync**: Synchronize data from external APIs (Cupid API integration)
- **Search & Indexing**: Optimized database queries with search indices
- **Rate Limiting**: Built-in request rate limiting
- **Database Migrations**: Version-controlled database schema management

## Technology Stack

- **Language**: Go 1.24.3
- **Database**: PostgreSQL
- **HTTP Router**: julienschmidt/httprouter
- **Scheduling**: go-co-op/gocron
- **Database Driver**: lib/pq (PostgreSQL)
- **Testing**: DATA-DOG/go-sqlmock
- **Rate Limiting**: golang.org/x/time

## Project Structure

```
nuitee/
├── cmd/
│   ├── api/          # REST API server
│   └── sync/         # Data synchronization service
├── internal/
│   ├── cache/        # Caching layer
│   ├── data/         # Data models and database operations
│   └── validator/    # Input validation
├── migrations/       # Database migrations
├── deployment/       # Deployment configurations
├── docs/            # Documentation
├── scripts/         # Utility scripts
└── bin/             # Compiled binaries
```

## Prerequisites

- Go 1.24.3 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

## Environment Variables

Create a .envrc file in the project root with your database configuration:

```env
export NUITEE_DB_DSN='postgres://username:password@localhost/nuitee?sslmode=disable'
export OPEN_AI_KEY='your_openai_api_key'
export CUPID_API_KEY='your_cupid_api_key'
export CUPID_API_URL='https://api.cupid.example.com'
export SYNC_INTERVAL=3
```

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/JLL32/nuitee.git
   cd nuitee
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up the database:**
   ```bash
   # Create PostgreSQL database
   createdb nuitee

   # Run migrations
   make db/migrations/up
   ```

## Usage

### Running the API Server

```bash
# Using Make
make run/api

# Or directly with Go
go run ./cmd/api -db-dsn="your_db_dsn" -openai-key="your_openai_key"
```

The API server will start on the configured port (default: 4000).

### Running Data Sync

```bash
# Using Make
make run/sync

# Or directly with Go
go run ./cmd/sync -db-dsn="your_db_dsn" -api-key="your_api_key" -api-url="api_url" -input="input.txt"
```

### Available Make Commands

#### Development
- `make run/api` - Run the API server
- `make run/sync` - Run the data synchronization

#### Database
- `make db/psql` - Connect to database using psql
- `make db/migrations/new name=migration_name` - Create new migration
- `make db/migrations/up` - Apply all migrations
- `make db/migrations/down` - Rollback last migration
- `make db/migrations/version` - Show current migration version
- `make db/reset` - Reset database (drop and recreate)

#### Testing
- `make test` - Run unit tests
- `make test/verbose` - Run tests with verbose output
- `make test/coverage` - Run tests with coverage report
- `make test/coverage/html` - Generate HTML coverage report
- `make test/race` - Run tests with race condition detection
- `make test/bench` - Run benchmark tests

#### Quality Control
- `make audit` - Run comprehensive quality checks
- `make vendor` - Tidy and vendor dependencies

#### Build
- `make build/api` - Build API binary for multiple platforms

## API Endpoints

The REST API provides the following endpoints:

### System Endpoints
- `GET /v1/healthcheck` - Health check endpoint
- `GET /debug/vars` - Runtime metrics and variables

### Documentation Endpoints
- `GET /docs` - API documentation index
- `GET /docs/swagger` - Swagger UI
- `GET /docs/redoc` - ReDoc UI
- `GET /docs/simple` - Simple HTML documentation
- `GET /docs/openapi.yaml` - OpenAPI specification

### Hotel Endpoints
- `GET /v1/hotels` - List hotels with filtering and pagination
- `GET /v1/hotels/:hotelID` - Get specific hotel details

### Review Endpoints
- `GET /v1/hotels/:hotelID/reviews` - Get reviews for a specific hotel
- `GET /v1/hotels/:hotelID/reviews/:reviewID` - Get specific review details
- `GET /v1/hotels/:hotelID/reviews/:reviewID/summary` - Get AI-generated review summary

All API endpoints are versioned with `/v1/` prefix and use RESTful conventions.

## Database Schema

### Hotels Table
- `hotel_id` - Primary key
- `hotel_name` - Hotel name
- `address`, `city`, `state`, `country`, `postal_code` - Location details
- `stars` - Star rating
- `rating` - Average rating
- `review_count` - Number of reviews
- `child_allowed`, `pets_allowed` - Amenity flags
- `description` - Hotel description

### Reviews Table
- `id` - Primary key
- `hotel_id` - Foreign key to hotels
- `average_score` - Review score
- `headline` - Review headline
- `pros`, `cons` - Review content
- `source` - Review source
- `language` - Review language

## Testing

Run tests with coverage:

```bash
make test/coverage/html
```

This generates a detailed HTML coverage report at `coverage.html`.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run quality checks (`make audit`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For questions or support, please open an issue on GitHub or contact the maintainers.
