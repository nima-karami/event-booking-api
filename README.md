.PHONY: docker-up docker-down docker-logs docker-clean dev

# Start PostgreSQL container
docker-up:
    docker-compose up -d

# Stop PostgreSQL container
docker-down:
    docker-compose down

# View PostgreSQL logs
docker-logs:
    docker-compose logs -f postgres

# Stop and remove container + volume (deletes all data)
docker-clean:
    docker-compose down -v

# Start development server (requires docker-up first)
dev:
    air