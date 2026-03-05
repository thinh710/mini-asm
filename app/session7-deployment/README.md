# 🚢 Buổi 6: Frontend Integration & Deployment

## Mục Tiêu

- ✅ Simple frontend SPA với vanilla JavaScript
- ✅ CORS configuration
- ✅ Docker containerization
- ✅ Docker Compose for full stack
- ✅ API documentation với Swagger
- ✅ Deployment basics

## Full Stack Architecture

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP/JSON
┌──────▼──────┐
│  Frontend   │  (HTML/CSS/JS)
│  (Port 80)  │
└──────┬──────┘
       │ API calls
┌──────▼──────────┐
│   Go Backend    │
│   (Port 8080)   │
└──────┬──────────┘
       │ SQL
┌──────▼──────────┐
│   PostgreSQL    │
│   (Port 5432)   │
└─────────────────┘
```

## Project Structure

```
session6-deployment/
├── web/
│   ├── index.html
│   ├── app.js
│   └── styles.css
├── cmd/server/
│   └── main.go
├── internal/
│   └── ...
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── nginx.conf
└── README.md
```

## 1. Simple Frontend

### index.html

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Mini ASM Dashboard</title>
    <link rel="stylesheet" href="styles.css" />
  </head>
  <body>
    <div class="container">
      <h1>🛡️ Attack Surface Management</h1>

      <!-- Create Asset Form -->
      <div class="card">
        <h2>Add New Asset</h2>
        <form id="createForm">
          <input
            type="text"
            id="name"
            placeholder="Name (e.g., example.com)"
            required
          />
          <select id="type" required>
            <option value="">Select Type</option>
            <option value="domain">Domain</option>
            <option value="ip">IP Address</option>
            <option value="service">Service</option>
          </select>
          <button type="submit">Add Asset</button>
        </form>
      </div>

      <!-- Filters -->
      <div class="card">
        <h2>Filters</h2>
        <div class="filters">
          <select id="filterType">
            <option value="">All Types</option>
            <option value="domain">Domain</option>
            <option value="ip">IP</option>
            <option value="service">Service</option>
          </select>
          <select id="filterStatus">
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
          <input type="text" id="search" placeholder="Search..." />
        </div>
      </div>

      <!-- Assets List -->
      <div class="card">
        <h2>Assets</h2>
        <div id="assetsList"></div>
      </div>
    </div>

    <script src="app.js"></script>
  </body>
</html>
```

### app.js

```javascript
const API_URL = "http://localhost:8080";

// Load assets on page load
document.addEventListener("DOMContentLoaded", () => {
  loadAssets();
  setupEventListeners();
});

function setupEventListeners() {
  document.getElementById("createForm").addEventListener("submit", createAsset);
  document.getElementById("filterType").addEventListener("change", loadAssets);
  document
    .getElementById("filterStatus")
    .addEventListener("change", loadAssets);
  document
    .getElementById("search")
    .addEventListener("input", debounce(loadAssets, 300));
}

async function loadAssets() {
  try {
    const type = document.getElementById("filterType").value;
    const status = document.getElementById("filterStatus").value;
    const search = document.getElementById("search").value;

    let url = `${API_URL}/assets?`;
    if (type) url += `type=${type}&`;
    if (status) url += `status=${status}&`;
    if (search) url += `search=${search}`;

    const response = await fetch(url);
    if (!response.ok) throw new Error("Failed to fetch assets");

    const assets = await response.json();
    displayAssets(assets);
  } catch (error) {
    console.error("Error loading assets:", error);
    showError("Failed to load assets");
  }
}

async function createAsset(e) {
  e.preventDefault();

  const name = document.getElementById("name").value;
  const type = document.getElementById("type").value;

  try {
    const response = await fetch(`${API_URL}/assets`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name, type }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to create asset");
    }

    document.getElementById("createForm").reset();
    loadAssets();
    showSuccess("Asset created successfully");
  } catch (error) {
    console.error("Error creating asset:", error);
    showError(error.message);
  }
}

async function deleteAsset(id) {
  if (!confirm("Are you sure you want to delete this asset?")) return;

  try {
    const response = await fetch(`${API_URL}/assets/${id}`, {
      method: "DELETE",
    });

    if (!response.ok) throw new Error("Failed to delete asset");

    loadAssets();
    showSuccess("Asset deleted successfully");
  } catch (error) {
    console.error("Error deleting asset:", error);
    showError("Failed to delete asset");
  }
}

function displayAssets(assets) {
  const container = document.getElementById("assetsList");

  if (!assets || assets.length === 0) {
    container.innerHTML = '<p class="empty">No assets found</p>';
    return;
  }

  container.innerHTML = assets
    .map(
      (asset) => `
        <div class="asset-card">
            <div class="asset-header">
                <h3>${asset.name}</h3>
                <span class="badge badge-${asset.type}">${asset.type}</span>
                <span class="badge badge-${asset.status}">${asset.status}</span>
            </div>
            <div class="asset-meta">
                <span>Created: ${new Date(asset.created_at).toLocaleDateString()}</span>
                <button onclick="deleteAsset('${asset.id}')" class="btn-delete">Delete</button>
            </div>
        </div>
    `,
    )
    .join("");
}

function debounce(func, wait) {
  let timeout;
  return function (...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), wait);
  };
}

function showSuccess(message) {
  // Simple alert for now - can enhance with toast notifications
  alert(message);
}

function showError(message) {
  alert(`Error: ${message}`);
}
```

## 2. CORS Configuration

```go
// internal/middleware/cors.go
package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow requests from frontend
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Update main.go
func main() {
    // ... setup handlers

    // Wrap with CORS middleware
    handler := middleware.CORS(mux)

    http.ListenAndServe(":8080", handler)
}
```

## 3. Dockerfile for Backend

```dockerfile
# Multi-stage build for smaller image
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./server"]
```

## 4. Docker Compose - Full Stack

```yaml
version: "3.8"

services:
  # PostgreSQL Database
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
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Go Backend API
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=mini_asm
      - DB_SSLMODE=disable
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  # Nginx for Frontend
  frontend:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./web:/usr/share/nginx/html:ro
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - api
    restart: unless-stopped

volumes:
  pgdata:
```

## 5. Nginx Configuration

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    server {
        listen 80;
        server_name localhost;

        # Serve frontend
        location / {
            root /usr/share/nginx/html;
            index index.html;
            try_files $uri $uri/ /index.html;
        }

        # Proxy API requests
        location /api/ {
            proxy_pass http://api:8080/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }
}
```

## 6. Swagger Documentation

```yaml
# docs/swagger.yaml
openapi: 3.0.0
info:
  title: Mini ASM API
  version: 1.0.0
  description: Attack Surface Management API

servers:
  - url: http://localhost:8080
    description: Local development

paths:
  /assets:
    get:
      summary: List all assets
      parameters:
        - name: type
          in: query
          schema:
            type: string
            enum: [domain, ip, service]
        - name: status
          in: query
          schema:
            type: string
            enum: [active, inactive]
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Asset"
    # ... other endpoints

components:
  schemas:
    Asset:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        type:
          type: string
          enum: [domain, ip, service]
        status:
          type: string
          enum: [active, inactive]
        created_at:
          type: string
          format: date-time
```

## Deployment Commands

```bash
# Build and start all services
docker-compose up --build

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild specific service
docker-compose build api
docker-compose up -d api
```

## Testing Full Stack

1. **Start services:**

   ```bash
   docker-compose up -d
   ```

2. **Check health:**

   ```bash
   curl http://localhost:8080/health
   ```

3. **Open frontend:**

   ```
   http://localhost
   ```

4. **Test API directly:**
   ```bash
   curl http://localhost:8080/assets
   ```

## Teaching Flow

### 1. Frontend Demo (20 phút)

- Show HTML structure
- Explain Fetch API
- Test locally without backend

### 2. CORS Explanation (15 phút)

- Why CORS?
- Preflight requests
- Middleware implementation

### 3. Docker Basics (20 phút)

- Dockerfile walkthrough
- Multi-stage builds
- Image layers

### 4. Docker Compose (20 phút)

- Service dependencies
- Health checks
- Volumes

### 5. Full Stack Demo (30 phút)

- Build and start
- Test end-to-end
- Show logs
- Troubleshooting

### 6. Deployment Options (15 phút)

- Cloud platforms (AWS, GCP, Azure)
- Kubernetes basics
- CI/CD pipeline

## Production Considerations

- ✅ Environment variables for secrets
- ✅ HTTPS/TLS configuration
- ✅ Rate limiting
- ✅ Authentication & Authorization
- ✅ Monitoring & Logging (Prometheus, Grafana)
- ✅ Backup strategies
- ✅ Auto-scaling
- ✅ CDN for static assets

## Homework

1. **Add authentication** (JWT)
2. **Implement rate limiting**
3. **Setup CI/CD** with GitHub Actions
4. **Deploy to cloud** (Heroku/Railway/Render)
5. **Add monitoring** (Prometheus + Grafana)

## Resources

- [Docker Documentation](https://docs.docker.com/)
- [Nginx Configuration](https://nginx.org/en/docs/)
- [Swagger Editor](https://editor.swagger.io/)
- [12 Factor App](https://12factor.net/)
