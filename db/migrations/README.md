# Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema versioning.

## Migration Files Location

Migrations are stored in `db/migrations/` directory.

## Migration File Naming Convention

Files follow the pattern: `{version}_{description}.{up|down}.sql`

Example:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

## Available Commands

### Using Makefile (Recommended)

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Show current migration version
make migrate-version

# Create new migration
make migrate-create NAME=add_user_avatar

# Force specific version (use with caution!)
make migrate-force V=1
```

### Using Go Application

Migrations are automatically applied when the application starts via `db.InitDB()`.

## Manual Migration Management

### Install migrate CLI tool (optional)

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Manual Commands

```bash
# Up
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up

# Down
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" down

# Force version
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" force VERSION
```

## Creating New Migrations

1. **Using Makefile:**
   ```bash
   make migrate-create NAME=add_column_to_users
   ```

2. **Manually:** Create two files:
   - `db/migrations/000004_description.up.sql` - Changes to apply
   - `db/migrations/000004_description.down.sql` - How to revert

## Migration Best Practices

1. **Always create both up and down migrations** - Ensure rollback capability
2. **Test migrations locally first** - Verify on development database
3. **Keep migrations small and focused** - One logical change per migration
4. **Never modify existing migrations** - Create new ones instead
5. **Use transactions when possible** - Wrap multiple statements in BEGIN/COMMIT
6. **Add indexes for foreign keys** - Improve query performance
7. **Document complex migrations** - Add comments explaining the why

## Schema Versioning Table

Migrations are tracked in the `schema_migrations` table, which is automatically created and managed by golang-migrate.

## Troubleshooting

### Dirty Migration State

If a migration fails midway, the database may be in a "dirty" state:

```bash
# Check current version and dirty status
make migrate-version

# Force to a specific version (after manually fixing database)
make migrate-force V=3
```

### Reset Database (Development Only)

```bash
# Drop all tables and start fresh
make docker-clean
make docker-up
# Migrations will run automatically on next app start
```
