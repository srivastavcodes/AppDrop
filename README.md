# AppDrop Backend Assignment

A REST API for managing mobile app pages and widgets, built with Go and PostgreSQL.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool

Install golang-migrate:
```bash
# macOS
brew install golang-migrate

# Or with Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Quick Start

### 1. Clone and Install Dependencies

```bash
git clone https://github.com/srivastavcodes/AppDrop

cd AppDrop

go mod download
```

### 2. Database Setup

Create PostgreSQL database:
```bash
createdb appdrop
```

Configure database connection (create `.envrc` file):
```bash
export APP_DROP_DSN="postgres://username:password@localhost/appdrop?sslmode=disable"
```

Or copy the example:
```bash
cp .envrc .your_envrc
# Edit .envrc with your database credentials
```

Run migrations:
```bash
make db/mig/up
```

### 3. Run the Application

```bash
# Using Makefile
make run/api

# Or directly
source .envrc
go run ./cmd/api -db-dsn="${APP_DROP_DSN}"
```

Server starts on `http://localhost:4000`

## Project Structure

```
.
├── cmd/api/                   # Application entry point
│   ├── main.go                # Server setup
│   ├── routes.go              # Route definitions
│   ├── pages.go               # Page endpoints
│   ├── widgets.go             # Widget endpoints
│   ├── middleware.go          # HTTP middleware
│   ├── errors.go              # Error handling
│   ├── helpers.go             # Helper functions
│   └── helpers.go             # Utilities
├── internal/data/             # Data layer
│   ├── models.go              # Model aggregation
│   ├── pages.go               # Page operations
│   └── widgets.go             # Widget operations
├── migrations/                # Database migrations
│   ├── 000001_create_pages_table.up.sql
│   ├── 000001_create_pages_table.down.sql
│   ├── 000002_create_widgets_table.up.sql
│   └── 000002_create_widgets_table.down.sql
├── .envrc                     # Environment variables
├── Makefile                   # Build commands
└── go.mod                     # Dependencies
```

## API Endpoints

### Pages
- `GET /pages` - List all pages
- `GET /pages/:id` - Get page with widgets
- `POST /pages` - Create page
- `PUT /pages/:id` - Update page
- `DELETE /pages/:id` - Delete page

### Widgets
- `POST /pages/:id/widgets` - Add widget to page
- `PUT /widgets/:id` - Update widget
- `DELETE /widgets/:id` - Delete widget
- `POST /pages/:id/widgets/reorder` - Reorder widgets

## Example Requests

### Create a page
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home",
    "route": "/home",
    "is_home": true
  }'
```

### Add a widget
```bash
curl -X POST http://localhost:4000/pages/{PAGE_ID}/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "banner",
    "config": {
      "image_url": "https://example.com/banner.jpg",
      "link": "/collections/new"
    }
  }'
```

### Get page with widgets
```bash
curl http://localhost:4000/pages/{PAGE_ID}
```

**For complete API documentation, see:**
- `API_REQUESTS.md` - Comprehensive examples

**For automated testing:**
```bash
./test_api.sh
```

## Database Schema

### Pages Table
- `id` (UUID) - Primary key
- `app_id` (UUID) - Application identifier (to mimic relations/constraints)
- `name` (VARCHAR) - Page name
- `route` (VARCHAR) - Unique route per app
- `is_home` (BOOLEAN) - Only one per app
- `created_at`, `updated_at` (TIMESTAMP)

Constraints:
- Unique index on `(app_id, route)`
- Unique partial index on `app_id` where `is_home = true`

### Widgets Table
- `id` (UUID) - Primary key
- `page_id` (UUID) - Foreign key to pages (CASCADE delete)
- `type` (VARCHAR) - Widget type
- `position` (INTEGER) - Order on page
- `config` (JSONB) - Flexible configuration
- `created_at`, `updated_at` (TIMESTAMP)

Constraints:
- Valid types: `banner`, `product_grid`, `text`, `image`, `spacer`
- Indexes on `(page_id, position)` and `page_id`

## Validation Rules

**Pages:**
- Name is required and cannot be empty
- Route must be unique per app
- Only one page can have `is_home = true` per app
- Cannot delete the home page

**Widgets:**
- Type must be one of: `banner`, `product_grid`, `text`, `image`, `spacer`
- Config is optional but must be valid JSON if provided
- Position is auto-assigned on creation

## Error Response Format

All errors follow this structure:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "page route already exists for this app"
  }
}
```

Error codes: `VALIDATION_ERROR`, `CONFLICT`, `NOT_FOUND`, `BAD_REQUEST`, `SERVER_ERROR`

## Makefile Commands

```bash
make run/api          # Run the application
make db/mig/up        # Run migrations up
make db/mig/down      # Run migrations down
make db/mig/new name=migration_name  # Create new migration
```

## Dependencies

```go
go 1.25

require (
    github.com/google/uuid          // UUID generation
    github.com/julienschmidt/httprouter  // HTTP router
    github.com/lib/pq               // PostgreSQL driver
)
```

Install:
```bash
go get github.com/google/uuid
go get github.com/julienschmidt/httprouter
go get github.com/lib/pq
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_DROP_DSN` | PostgreSQL connection string | `postgres://user:pass@localhost/appdrop?sslmode=disable` |

