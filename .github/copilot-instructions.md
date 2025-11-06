# Event Booking API - AI Agent Instructions

## Architecture Overview

This is a Go REST API using Gin framework with SQLite database for event management and user registration.

**Key architectural decisions:**
- **No ORM**: Raw SQL queries with `database/sql` package directly (see `models/` for pattern)
- **Layered structure**: Routes → Models → Database, with middleware for cross-cutting concerns
- **Shared DB connection**: Global `db.DB` variable initialized in `main.go`, reused across all models
- **Model methods**: Each model has CRUD methods (`Save()`, `Update()`, `Delete()`) that operate on receiver instances

## Project Structure

```
main.go              # Entry point: InitDB → RegisterRoutes → server.Run(":8080")
db/db.go            # Database initialization, creates 3 tables: users, events, registrations
models/             # Business logic with direct SQL queries
├── event.go        # Event CRUD + Register/Unregister methods
├── user.go         # User CRUD + Authenticate method
└── registration.go # Join queries for event registrations
routes/             # HTTP handlers organized by resource
├── routes.go       # Central route registration with auth middleware
├── events.go       # Event endpoints
└── users.go        # User/auth endpoints
middlewares/        # Cross-cutting concerns
└── auth.go         # JWT validation, extracts userID to context
utils/              # Helpers
├── jwt.go          # Token generation/verification (hardcoded secret)
└── hash.go         # bcrypt password hashing
```

## Development Workflow

**Running the app:**
- Use Air for hot reload: `air` (config in `.air.toml`)
- Builds to `tmp/main.exe`, runs on port 8080
- Database file: `api.db` (auto-created on first run)

**No tests currently exist** - test files would need to be created if adding test coverage.

## Authentication Pattern

**Two route groups:**
1. **Public routes**: Directly on server (e.g., `server.POST("/users/signup")`)
2. **Protected routes**: Under `authenticated` group with middleware (e.g., `authenticated.POST("/events")`)

**Auth flow:**
1. Middleware checks `Authorization: Bearer <token>` header (`middlewares/auth.go`)
2. Validates JWT using `utils.VerifyToken()`
3. Sets `userID` in Gin context: `c.Set("userID", userID)`
4. Handlers retrieve with `c.GetInt64("userID")`

## Key Patterns & Conventions

### Model Pattern
Models define both struct and database operations:
```go
type Event struct {
    ID int64 `json:"id"`
    Title string `json:"title" binding:"required"`
    // ...
}

func (e *Event) Save() error {
    result, err := db.DB.Exec("INSERT INTO events (...) VALUES (...)", ...)
    e.ID, _ = result.LastInsertId()  // Mutates receiver with generated ID
    return err
}
```

### Error Responses
Consistent JSON structure across all handlers:
```go
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
```

### Authorization Checks
Owner-based authorization for updates/deletes:
```go
event, _ := models.GetEventByID(eventId)
if event.UserID != userID {
    c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized..."})
    return
}
```

### ID Parsing
Convert URL params to int64:
```go
eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
```

## Database Schema

- **users**: id, email (unique), password (bcrypt hashed)
- **events**: id, title, description, location, date, user_id (FK to users)
- **registrations**: id, user_id (FK), event_id (FK) - many-to-many join table

Tables auto-created via `db.createTables()` on startup using `CREATE TABLE IF NOT EXISTS`.

## Dependencies

- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/ncruces/go-sqlite3` - SQLite driver (pure Go implementation)
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `golang.org/x/crypto/bcrypt` - Password hashing

## Important Notes

- JWT secret is hardcoded as `"your_secret_key"` in `utils/jwt.go` - should be environment variable
- Token expiration: 2 hours
- Bcrypt cost: 14
- Max DB connections: 10 open, 5 idle
- No migration system - schema changes require manual SQL or dropping tables
