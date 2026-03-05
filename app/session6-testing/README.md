# 🧪 Buổi 5: Testing & Quality Assurance

## Mục Tiêu

- ✅ Viết unit tests cho từng layer
- ✅ Viết integration tests cho API
- ✅ Mock dependencies với interfaces
- ✅ Test coverage measurement
- ✅ Table-driven tests pattern
- ✅ Logging và monitoring
- ✅ Benchmarking

## Testing Philosophy

> "If it's not tested, it's broken"

```
Testing Pyramid:
       /\
      /E2\     ← Few (End-to-end)
     /----\
    /Integ\   ← Some (Integration)
   /--------\
  /   Unit   \ ← Many (Unit tests)
 /____________\
```

## Test Structure

```
session5-testing/
├── internal/
│   ├── model/
│   │   ├── asset.go
│   │   └── asset_test.go
│   ├── service/
│   │   ├── asset_service.go
│   │   └── asset_service_test.go
│   ├── handler/
│   │   ├── asset_handler.go
│   │   └── asset_handler_test.go
│   └── storage/
│       ├── storage.go
│       ├── memory/
│       │   ├── memory.go
│       │   └── memory_test.go
│       └── mock/
│           └── mock_storage.go
├── test/
│   ├── integration_test.go
│   └── benchmark_test.go
└── README.md
```

## 1. Unit Tests - Model Layer

```go
// internal/model/asset_test.go
package model

import "testing"

func TestIsValidType(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid domain", "domain", true},
        {"valid ip", "ip", true},
        {"valid service", "service", true},
        {"invalid type", "invalid", false},
        {"empty string", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidType(tt.input)
            if got != tt.expected {
                t.Errorf("IsValidType(%q) = %v, want %v",
                    tt.input, got, tt.expected)
            }
        })
    }
}
```

## 2. Unit Tests - Service Layer with Mocking

```go
// internal/storage/mock/mock_storage.go
package mock

import (
    "mini-asm/internal/model"
)

type MockStorage struct {
    CreateFunc   func(*model.Asset) error
    GetAllFunc   func() ([]*model.Asset, error)
    GetByIDFunc  func(string) (*model.Asset, error)
    // ... other methods
}

func (m *MockStorage) Create(asset *model.Asset) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(asset)
    }
    return nil
}

// ... implement other methods

// internal/service/asset_service_test.go
package service

import (
    "testing"
    "mini-asm/internal/model"
    "mini-asm/internal/storage/mock"
)

func TestAssetService_CreateAsset(t *testing.T) {
    tests := []struct {
        name      string
        inputName string
        inputType string
        wantErr   bool
        errType   error
    }{
        {
            name:      "valid asset",
            inputName: "example.com",
            inputType: "domain",
            wantErr:   false,
        },
        {
            name:      "empty name",
            inputName: "",
            inputType: "domain",
            wantErr:   true,
            errType:   model.ErrEmptyName,
        },
        {
            name:      "invalid type",
            inputName: "test",
            inputType: "invalid",
            wantErr:   true,
            errType:   model.ErrInvalidType,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockStorage := &mock.MockStorage{
                CreateFunc: func(asset *model.Asset) error {
                    return nil
                },
            }
            service := NewAssetService(mockStorage)

            // Act
            asset, err := service.CreateAsset(tt.inputName, tt.inputType)

            // Assert
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error but got nil")
                }
                if tt.errType != nil && err != tt.errType {
                    t.Errorf("expected error %v, got %v", tt.errType, err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if asset == nil {
                    t.Error("expected asset but got nil")
                }
                if asset != nil {
                    if asset.ID == "" {
                        t.Error("expected asset ID to be generated")
                    }
                    if asset.Status != model.StatusActive {
                        t.Errorf("expected status %s, got %s",
                            model.StatusActive, asset.Status)
                    }
                }
            }
        })
    }
}
```

## 3. Integration Tests - API Level

```go
// test/integration_test.go
package test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "mini-asm/internal/handler"
    "mini-asm/internal/model"
    "mini-asm/internal/service"
    "mini-asm/internal/storage/memory"
)

func setupTestServer() *http.ServeMux {
    store := memory.NewMemoryStorage()
    svc := service.NewAssetService(store)
    h := handler.NewAssetHandler(svc)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /assets", h.CreateAsset)
    mux.HandleFunc("GET /assets", h.ListAssets)
    mux.HandleFunc("GET /assets/{id}", h.GetAsset)

    return mux
}

func TestCreateAsset_Integration(t *testing.T) {
    server := setupTestServer()

    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        checkResponse  func(*testing.T, []byte)
    }{
        {
            name:           "valid asset",
            requestBody:    `{"name":"example.com","type":"domain"}`,
            expectedStatus: http.StatusCreated,
            checkResponse: func(t *testing.T, body []byte) {
                var asset model.Asset
                if err := json.Unmarshal(body, &asset); err != nil {
                    t.Fatalf("failed to unmarshal response: %v", err)
                }
                if asset.ID == "" {
                    t.Error("expected ID to be set")
                }
                if asset.Name != "example.com" {
                    t.Errorf("expected name example.com, got %s", asset.Name)
                }
            },
        },
        {
            name:           "empty name",
            requestBody:    `{"name":"","type":"domain"}`,
            expectedStatus: http.StatusBadRequest,
            checkResponse: func(t *testing.T, body []byte) {
                var errResp map[string]string
                json.Unmarshal(body, &errResp)
                if errResp["error"] == "" {
                    t.Error("expected error message")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create request
            req := httptest.NewRequest("POST", "/assets",
                bytes.NewBufferString(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")

            // Record response
            w := httptest.NewRecorder()
            server.ServeHTTP(w, req)

            // Check status code
            if w.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d",
                    tt.expectedStatus, w.Code)
            }

            // Check response body
            if tt.checkResponse != nil {
                tt.checkResponse(t, w.Body.Bytes())
            }
        })
    }
}

func TestFullCRUDWorkflow_Integration(t *testing.T) {
    server := setupTestServer()

    // 1. Create asset
    createReq := httptest.NewRequest("POST", "/assets",
        bytes.NewBufferString(`{"name":"test.com","type":"domain"}`))
    createReq.Header.Set("Content-Type", "application/json")
    createW := httptest.NewRecorder()
    server.ServeHTTP(createW, createReq)

    if createW.Code != http.StatusCreated {
        t.Fatalf("create failed: %d", createW.Code)
    }

    var createdAsset model.Asset
    json.Unmarshal(createW.Body.Bytes(), &createdAsset)

    // 2. Get asset
    getReq := httptest.NewRequest("GET", "/assets/"+createdAsset.ID, nil)
    getW := httptest.NewRecorder()
    server.ServeHTTP(getW, getReq)

    if getW.Code != http.StatusOK {
        t.Errorf("get failed: %d", getW.Code)
    }

    // 3. List assets
    listReq := httptest.NewRequest("GET", "/assets", nil)
    listW := httptest.NewRecorder()
    server.ServeHTTP(listW, listReq)

    if listW.Code != http.StatusOK {
        t.Errorf("list failed: %d", listW.Code)
    }

    var assets []*model.Asset
    json.Unmarshal(listW.Body.Bytes(), &assets)
    if len(assets) != 1 {
        t.Errorf("expected 1 asset, got %d", len(assets))
    }
}
```

## 4. Benchmark Tests

```go
// test/benchmark_test.go
package test

import (
    "testing"
    "mini-asm/internal/service"
    "mini-asm/internal/storage/memory"
)

func BenchmarkCreateAsset(b *testing.B) {
    store := memory.NewMemoryStorage()
    svc := service.NewAssetService(store)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        svc.CreateAsset("example.com", "domain")
    }
}

func BenchmarkGetAllAssets(b *testing.B) {
    store := memory.NewMemoryStorage()
    svc := service.NewAssetService(store)

    // Seed with 1000 assets
    for i := 0; i < 1000; i++ {
        svc.CreateAsset(fmt.Sprintf("example%d.com", i), "domain")
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        svc.GetAllAssets()
    }
}
```

## 5. Logging

```go
// pkg/logger/logger.go
package logger

import (
    "log"
    "os"
)

var (
    InfoLogger  *log.Logger
    ErrorLogger *log.Logger
)

func Init() {
    InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Usage in handlers
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
    logger.InfoLogger.Printf("Creating asset: %s", req.Name)

    asset, err := h.service.CreateAsset(req.Name, req.Type)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to create asset: %v", err)
        // ...
    }
}
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test -run TestCreateAsset ./internal/service

# Run benchmarks
go test -bench=. ./test

# Run with race detector
go test -race ./...

# Verbose output
go test -v ./...
```

## Coverage Goals

- **Service layer:** > 80% coverage
- **Handler layer:** > 70% coverage
- **Storage layer:** > 60% coverage
- **Overall:** > 70% coverage

## Teaching Flow

### 1. Introduction to Testing (15 phút)

- Why test?
- Types of tests
- Testing pyramid

### 2. Unit Tests - Model (20 phút)

- Table-driven tests
- Test naming conventions
- Assertions

### 3. Mocking & Service Tests (30 phút)

- Why mock?
- Create mock storage
- Test service in isolation

### 4. Integration Tests (30 phút)

- httptest package
- Full request/response cycle
- Test workflows

### 5. Coverage & Benchmarks (20 phút)

- Generate coverage report
- Identify untested code
- Performance benchmarks

### 6. Best Practices (15 phút)

- Test organization
- Common mistakes
- CI/CD integration

## Homework

1. **Achieve 80%+ test coverage**
2. **Add tests for error cases**
3. **Write benchmark comparisons** (memory vs database)
4. **Add test fixtures** (sample data)
5. **Setup GitHub Actions** for CI

## Resources

- [Go Testing](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Coverage](https://blog.golang.org/cover)
