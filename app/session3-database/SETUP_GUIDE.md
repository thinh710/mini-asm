# 📋 Session 3 Setup & Testing Guide

## ✅ Pre-Flight Checklist

Before teaching this session, verify everything works:

### 1. Docker is Running

```powershell
docker --version
docker-compose --version
```

Expected: Docker version 20+ and Docker Compose version 2+

### 2. Go is Installed

```powershell
go version
```

Expected: Go 1.22 or higher

### 3. All Files Present

```
session3-database/
├── .env.example           ✓
├── .gitignore             ✓
├── docker-compose.yml     ✓
├── go.mod                 ✓
├── go.sum                 ✓
├── Makefile               ✓
├── setup.ps1              ✓
├── README.md              ✓
├── cmd/
│   └── server/
│       └── main.go        ✓
├── internal/
│   ├── handler/
│   │   ├── asset_handler.go  ✓
│   │   └── health_handler.go ✓
│   ├── model/
│   │   ├── asset.go       ✓
│   │   └── errors.go      ✓
│   ├── service/
│   │   └── asset_service.go  ✓
│   └── storage/
│       ├── storage.go     ✓
│       ├── memory/
│       │   └── memory.go  ✓
│       └── postgres/
│           └── postgres.go   ✓
└── migrations/
    ├── 001_create_assets.up.sql    ✓
    └── 001_create_assets.down.sql  ✓
```

---

## 🚀 Quick Test (5 minutes)

### Option A: Using PowerShell Script (Windows)

```powershell
# 1. Setup (starts DB and waits)
.\setup.ps1 setup

# 2. Start server (in this terminal)
.\setup.ps1 run
```

### Option B: Using Manual Commands

```powershell
# 1. Start database
docker-compose up -d

# 2. Wait for database to be ready
Start-Sleep -Seconds 5

# 3. Verify database is running
docker-compose ps

# 4. Start server
go run cmd/server/main.go
```

---

## 🧪 Testing Workflow

Open a **new terminal** and run these tests:

### Test 1: Health Check

```powershell
curl http://localhost:8080/health
```

Expected:

```json
{ "status": "ok", "timestamp": "2026-03-05T..." }
```

### Test 2: Create Asset

```powershell
curl -X POST http://localhost:8080/assets `
  -H "Content-Type: application/json" `
  -d '{\"name\":\"example.com\",\"type\":\"domain\"}'
```

Expected: Status 201 with asset object including `id`

### Test 3: List Assets

```powershell
curl http://localhost:8080/assets
```

Expected: Array with the asset created above

### Test 4: Get Single Asset

```powershell
# Copy the ID from previous response
curl http://localhost:8080/assets/{id}
```

Expected: Single asset object

### Test 5: **CRITICAL TEST - Persistence!**

```powershell
# 1. Stop the server (Ctrl+C in server terminal)

# 2. Restart the server
go run cmd/server/main.go

# 3. List assets again
curl http://localhost:8080/assets
```

Expected: **Assets still there!** This proves database persistence.

### Test 6: Verify in Database

```powershell
docker-compose exec db psql -U postgres -d mini_asm -c "SELECT * FROM assets;"
```

Expected: Table output showing all assets

---

### Q&A and Homework

Common questions:

- "Why PostgreSQL and not MySQL/MongoDB?"
- "What if database is down?"
- "How to backup data?"
- "What about SQL injection?"

Homework:

- Add database health check to `/health` endpoint
- Implement connection retry logic
- Add more indexes and test performance

---

## 🐛 Common Issues & Solutions

### Issue 1: Port 5432 Already in Use

**Symptom:**

```
Error starting userland proxy: listen tcp4 0.0.0.0:5432: bind: address already in use
```

**Solution:**

```powershell
# Find process using port 5432
Get-NetTCPConnection -LocalPort 5432

# Stop existing PostgreSQL service
Stop-Service postgresql-x64-15

# Or change port in docker-compose.yml
ports:
  - "5433:5432"  # Changed host port
```

### Issue 2: Migrations Not Applied

**Symptom:**

```
relation "assets" does not exist
```

**Solution:**

```powershell
# Check if migrations ran
.\setup.ps1 db-shell
\dt  # Should show 'assets' table

# If not, manually apply
.\setup.ps1 migrate-up
```

### Issue 3: Connection Refused

**Symptom:**

```
connection refused
```

**Check:**

```powershell
# 1. Is Docker running?
docker ps

# 2. Is database container running?
docker-compose ps

# 3. Check database logs
.\setup.ps1 db-logs

# 4. Restart database
.\setup.ps1 db-stop
.\setup.ps1 db-start
```

### Issue 4: Wrong Database Credentials

**Symptom:**

```
password authentication failed for user "postgres"
```

**Solution:**

- Check `.env` matches `docker-compose.yml`
- Default: user=postgres, password=postgres
- Reset database: `.\setup.ps1 clean` then `.\setup.ps1 setup`

---

## 📚 Additional Resources

### Official Documentation

- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [Go database/sql package](https://pkg.go.dev/database/sql)
- [Docker Compose docs](https://docs.docker.com/compose/)

### Recommended Reading

- Database design best practices
- SQL injection and how to prevent it
- Connection pooling in Go
- Database indexing strategies
