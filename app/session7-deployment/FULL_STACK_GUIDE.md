# EASM Platform - Complete Full-Stack Demo

Complete External Attack Surface Management platform with React frontend and Go backend.

## 🎯 Project Overview

```
session6-testing/
├── backend/                    # Go API Server
│   ├── cmd/server/main.go     # Server entry point
│   ├── internal/              # Application code
│   │   ├── handler/           # HTTP handlers
│   │   ├── service/           # Business logic
│   │   ├── storage/           # Data layer
│   │   ├── model/             # Domain models
│   │   └── validator/         # Input validation
│   └── migrations/            # Database migrations
│
└── frontend/                  # React SPA
    ├── src/
    │   ├── pages/            # Page components
    │   ├── services/         # API client
    │   └── App.jsx           # Main app
    └── package.json
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Node.js 18+
- PostgreSQL 13+
- Docker & Docker Compose (optional)

### Option 1: Using Docker Compose (Recommended)

```bash
# Start everything (database + backend)
docker-compose up -d

# Install frontend dependencies
cd frontend
npm install

# Start frontend dev server
npm run dev
```

**Access:**

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

### Option 2: Manual Setup

#### 1. Start Database

```bash
# Using Docker
docker-compose up -d postgres

# Or use your local PostgreSQL
# Create database: easm_db
```

#### 2. Configure Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your database credentials
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=easm_db
```

#### 3. Run Migrations

```bash
# Using migrate CLI
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/easm_db?sslmode=disable" up

# Or manually run SQL files in migrations/
```

#### 4. Start Backend

```bash
# Install dependencies
go mod download

# Run server
go run cmd/server/main.go

# Server starts on :8080
```

#### 5. Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Frontend starts on :3000
```

## 📺 Demo Walkthrough

### Step 1: Create Assets

1. Navigate to http://localhost:3000/assets
2. Click "Add Asset"
3. Add test domains:
   - `example.com` (type: domain)
   - `google.com` (type: domain)
   - `192.168.1.1` (type: ip)

### Step 2: Run Scans

1. Go to http://localhost:3000/scanning
2. Select an asset
3. Choose scan type:
   - **DNS** - Safe, queries public DNS
   - **WHOIS** - Domain registration info
   - **Subdomain** - Enumerate subdomains
4. Click "Start Scan"
5. Watch real-time status updates

### Step 3: View Results

1. Go to http://localhost:3000/results
2. Select your asset
3. Explore:
   - **DNS Records** - A, AAAA, MX, NS, TXT records
   - **Subdomains** - Discovered subdomains
   - **WHOIS** - Registration details

## 🎨 Frontend Features

### Dashboard

- System health status
- Asset statistics
- Quick start guide
- Feature overview

### Asset Management

- CRUD operations (Create, Read, Update, Delete)
- Filtering (by type, status)
- Search functionality
- Pagination
- Inline editing

### Scanning Operations

- Select asset and scan type
- Start scans with warnings for active scans
- Real-time status updates (auto-refresh)
- Scan history with details

### Results Visualization

- View all results or filter by type
- DNS records table
- Subdomains list with status
- WHOIS information display
- Raw data viewer

## 🔌 API Endpoints

**See [api.yml](api.yml) for complete OpenAPI specification**

### Health

- `GET /health` - System health check

### Assets

- `POST /assets` - Create asset
- `GET /assets` - List assets (with pagination, filters)
- `GET /assets/{id}` - Get single asset
- `PUT /assets/{id}` - Update asset
- `DELETE /assets/{id}` - Delete asset

### Scanning

- `POST /assets/{id}/scan` - Start scan
- `GET /assets/{id}/scans` - List scan jobs
- `GET /scan-jobs/{id}` - Get scan status
- `GET /scan-jobs/{id}/results` - Get results

### Results

- `GET /assets/{id}/results` - All results
- `GET /assets/{id}/dns` - DNS records
- `GET /assets/{id}/subdomains` - Subdomains
- `GET /assets/{id}/whois` - WHOIS data

## 🧪 Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Test Coverage:**

- **Models:** 76% (asset.go, scan.go)
- **Validators:** 74% (asset_validator.go)
- Total: 26 tests + 6 benchmarks

### Frontend Testing

```bash
cd frontend

# Run linter
npm run lint

# Build for production (validates code)
npm run build
```

## 📊 Architecture

### Backend Stack

- **Language:** Go 1.23
- **Framework:** net/http (stdlib)
- **Database:** PostgreSQL with pgx driver
- **Architecture:** Clean Architecture (3-layer)
  - Handler → Service → Storage

### Frontend Stack

- **Framework:** React 18
- **Build Tool:** Vite
- **HTTP Client:** Axios
- **Icons:** Lucide React
- **Routing:** React Router v6
- **Styling:** Custom CSS with variables

### Communication

```
[Browser] <--HTTP--> [Vite Dev Server] <--Proxy--> [Go API] <---> [PostgreSQL]
  :3000                    :3000                    :8080         :5432
```

## 🔒 Security Features

### Input Validation

- ✅ SQL injection prevention
- ✅ Null byte detection
- ✅ Domain format validation
- ✅ IP address validation
- ✅ Type checking

### Active Scan Warnings

- Frontend displays warnings
- Confirmation required
- Legal notices visible
- Only passive scans enabled by default

## 🎓 Educational Value

This project demonstrates:

1. **Full-Stack Development**
   - React frontend
   - Go backend
   - PostgreSQL database

2. **Clean Architecture**
   - Separation of concerns
   - Dependency injection
   - Interface-based design

3. **RESTful API Design**
   - Resource-based URLs
   - Proper HTTP methods
   - Standard status codes

4. **Real-time Updates**
   - Polling mechanism
   - Status tracking
   - Async operations

5. **Testing Best Practices**
   - Unit tests
   - Table-driven tests
   - Coverage analysis

## 🐛 Troubleshooting

### Backend won't start

```bash
# Check database connection
psql -h localhost -U postgres -d easm_db

# Check port availability
netstat -an | grep 8080

# Check logs
go run cmd/server/main.go
```

### Frontend can't connect

```bash
# Verify backend is running
curl http://localhost:8080/health

# Check proxy configuration in vite.config.js
# Clear cache
rm -rf node_modules/.vite
```

### Database errors

```bash
# Reset database
docker-compose down -v
docker-compose up -d

# Rerun migrations
migrate -path migrations -database "..." up
```

## 📚 Documentation

- **API Spec:** [api.yml](api.yml) - OpenAPI 3.0 specification
- **Frontend README:** [frontend/README.md](frontend/README.md)
- **Testing Guide:** [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **Architecture:** [SCANNING_ARCHITECTURE.md](SCANNING_ARCHITECTURE.md)

## 🚀 Deployment

### Frontend Build

```bash
cd frontend
npm run build

# Output: frontend/dist/
# Deploy to: Netlify, Vercel, S3, etc.
```

### Backend Build

```bash
# Build binary
go build -o easm-server cmd/server/main.go

# Run
./easm-server

# Or with Docker
docker build -t easm-backend .
docker run -p 8080:8080 easm-backend
```

## 🔜 Future Enhancements

- [ ] WebSocket for real-time updates
- [ ] GraphQL API option
- [ ] Export results (CSV, JSON, PDF)
- [ ] Scheduled scans (cron jobs)
- [ ] Email notifications
- [ ] Multi-user support with authentication
- [ ] RBAC (Role-Based Access Control)
- [ ] Dark mode UI
- [ ] Advanced filtering & search
- [ ] Result comparison over time
- [ ] Scan templates

## 📖 Learning Resources

- **Go:** https://go.dev/doc/tutorial
- **React:** https://react.dev/learn
- **PostgreSQL:** https://www.postgresql.org/docs/
- **Testing in Go:** https://go.dev/blog/subtests
- **Clean Architecture:** Martin Fowler's blog

## 🤝 Contributing

This is an educational project for CMC Intern Program - Session 6: Testing & Quality Assurance.

## 📝 License

Educational project - CMC Intern Program 2026

---

## ⭐ Demo Screenshots

### Dashboard

- System health status
- Asset statistics
- Feature overview

### Asset Management

- List view with filtering
- Create/Edit modal
- Inline actions

### Scanning

- Scan configuration
- Real-time job tracking
- Status updates

### Results

- DNS records table
- Subdomains list
- WHOIS viewer

---

**Built with ❤️ for learning purposes**

For questions, refer to documentation or contact your instructor.
