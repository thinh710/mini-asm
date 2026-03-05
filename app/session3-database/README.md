# 🗄️ Buổi 3: Database Integration

## ⚡ Quick Start

```bash
# 1. Start PostgreSQL database
docker-compose up -d

# 2. Verify database is running
docker-compose ps

# 3. Install dependencies
go mod tidy

# 4. Start the server
go run cmd/server/main.go

# 5. Test the API (in another terminal)
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"example.com","type":"domain"}'

curl http://localhost:8080/assets
```

**🎉 The Key Difference:** Restart the server - your data persists! Compare this with Session 2 where restarting lost all data.

---

## Mục Tiêu

- ✅ Thiết kế database schema
- ✅ Integration với PostgreSQL
- ✅ Database migration
- ✅ So sánh in-memory vs database persistence
- ✅ Configuration management

## So Sánh với Buổi 2

| Aspect           | Buổi 2 (Memory) | Buổi 3 (Database)              |
| ---------------- | --------------- | ------------------------------ |
| **Storage**      | In-memory map   | PostgreSQL table               |
| **Persistence**  | Lost on restart | Permanent                      |
| **Scalability**  | Single instance | Multiple instances             |
| **Code changes** | N/A             | **Only 1-2 lines in main.go!** |
| **Setup**        | None            | Docker + migrations            |

## Key Point: Clean Architecture Benefits

```go
// Buổi 2
store := memory.NewMemoryStorage()

// Buổi 3 - CHỈ THAY 1 DÒNG!
store := postgres.NewPostgresStorage(db)

// Handler, Service, Model: KHÔNG THAY ĐỔI!
```

## New Files

```
session3-database/
├── internal/
│   └── storage/
│       └── postgres/
│           └── postgres.go      # PostgreSQL implementation
├── migrations/
│   ├── 001_create_assets.up.sql
│   └── 001_create_assets.down.sql
├── .env.example
├── docker-compose.yml
└── README.md
```

## Database Schema

```sql
-- migrations/001_create_assets.up.sql
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('domain', 'ip', 'service')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_name ON assets(name);
CREATE INDEX idx_assets_created_at ON assets(created_at DESC);
```

## PostgreSQL Storage Implementation

Key differences from MemoryStorage:

1. **SQL Queries** instead of map operations
2. **Connection pooling** for performance
3. **Prepared statements** for security
4. **Transaction support** (future)

```go
// internal/storage/postgres/postgres.go
type PostgresStorage struct {
    db *sql.DB
}

func (p *PostgresStorage) Create(asset *model.Asset) error {
    query := `
        INSERT INTO assets (id, name, type, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := p.db.Exec(query,
        asset.ID,
        asset.Name,
        asset.Type,
        asset.Status,
        asset.CreatedAt,
        asset.UpdatedAt,
    )
    return err
}
```

## Docker Setup

```yaml
# docker-compose.yml
version: "3.8"
services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: mini_asm
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  pgdata:
```

## Configuration

```bash
# .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mini_asm
DB_SSLMODE=disable
```

## Flow

### Step 1: Review Buổi 2 (10 phút)

- Show memory storage working
- Restart server → data lost!
- "We need persistence!"

### Step 2: Database Design (20 phút)

- Draw ER diagram
- Explain schema design
- Column types and constraints
- Indexes for performance

### Step 3: Setup PostgreSQL (15 phút)

- Docker Compose walkthrough
- Start database: `docker-compose up -d`
- Connect with client: `psql -h localhost -U postgres mini_asm`
- Run migration

### Step 4: PostgresStorage Implementation (45 phút)

- Compare with MemoryStorage interface
- SQL queries walkthrough
- Prepared statements (security)
- Error handling (SQL errors → domain errors)

### Step 5: Configuration (15 phút)

- Environment variables
- Config loading
- Connection string building
- Connection pooling settings

### Step 6: Update Main (10 phút)

- Add database connection
- Swap storage implementation
- Test endpoints

### Step 7: Demo (15 phút)

1. Create assets via API
2. Stop server
3. Restart server
4. List assets → still there!
5. Show in database: `SELECT * FROM assets;`

### Step 8: Performance Comparison (10 phút)

- Benchmark memory vs database
- Explain trade-offs
- When to use each

## Key Concepts

### 1. SQL Injection Prevention

```go
// ❌ BAD - SQL injection vulnerable
query := fmt.Sprintf("SELECT * FROM assets WHERE id = '%s'", id)

// ✅ GOOD - Prepared statement
query := "SELECT * FROM assets WHERE id = $1"
db.Query(query, id)
```

### 2. Connection Pooling

```go
db.SetMaxOpenConns(25)      // Max connections
db.SetMaxIdleConns(5)       // Keep alive
db.SetConnMaxLifetime(5 * time.Minute)
```

### 3. Error Handling

```go
err := db.QueryRow(query, id).Scan(&asset)
if err == sql.ErrNoRows {
    return nil, model.ErrNotFound // Map to domain error
}
```

### 4. Migrations

- Up migration: create schema
- Down migration: rollback
- Version control for database schema
- Apply in order: 001, 002, 003...

## Testing

```bash
# Start database
docker-compose up -d

# Check database is running
docker ps

# Run migrations
psql -h localhost -U postgres -d mini_asm -f migrations/001_create_assets.up.sql

# Start server
go run cmd/server/main.go

# Test API
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"example.com","type":"domain"}'

# Verify in database
docker exec -it <container_id> psql -U postgres -d mini_asm -c "SELECT * FROM assets;"

# Restart server and verify data persists
```

## Homework

1. **Add database health check** to `/health` endpoint
   - Check if connection is alive
   - Return connection stats

2. **Implement connection retry logic**
   - Retry on connection failure
   - Exponential backoff

3. **Add database migrations tool**
   - Use golang-migrate or similar
   - Support up/down migrations

4. **Performance tuning**
   - Adjust connection pool settings
   - Add more indexes
   - Benchmark queries

## Common Issues

### Issue 1: Connection Refused

**Symptom:** `connection refused`
**Fix:** Check Docker is running, port 5432 not in use

### Issue 2: Password Authentication Failed

**Symptom:** `password authentication failed`
**Fix:** Check .env values match docker-compose.yml

### Issue 3: Table Already Exists

**Symptom:** `relation "assets" already exists`
**Fix:** Run down migration first, or drop database

## Resources

- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [Go database/sql](https://pkg.go.dev/database/sql)
- [lib/pq Driver](https://github.com/lib/pq)
- [Docker Compose Docs](https://docs.docker.com/compose/)
