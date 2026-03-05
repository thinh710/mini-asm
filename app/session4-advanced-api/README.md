# 🔍 Buổi 4: Advanced API Features

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

# 5. Test with advanced queries (in another terminal)
# Pagination
curl "http://localhost:8080/assets?page=1&page_size=10"

# Filtering
curl "http://localhost:8080/assets?type=domain&status=active"

# Search
curl "http://localhost:8080/assets?search=example"

# Sorting
curl "http://localhost:8080/assets?sort_by=name&sort_order=asc"

# Combined (everything together!)
curl "http://localhost:8080/assets?type=domain&search=google&page=2&sort_by=name&sort_order=asc"
```

---

## Mục Tiêu

- ✅ **Pagination**: Handle large datasets efficiently
- ✅ **Advanced Filtering**: Multiple filters combined
- ✅ **Search**: Partial name matching
- ✅ **Sorting**: Flexible sorting (asc/desc, multiple fields)
- ✅ **Input Validation**: Comprehensive type-specific validation
- ✅ **Error Handling**: Clear, actionable error messages

## So Sánh với Session 3

| Feature          | Session 3                   | Session 4 (Advanced)                  |
| ---------------- | --------------------------- | ------------------------------------- |
| **Pagination**   | ❌ None                     | ✅ ?page=1&page_size=20               |
| **Filtering**    | ✅ Basic (type, status)     | ✅ Enhanced + Search                  |
| **Sorting**      | ❌ Fixed order              | ✅ ?sort_by=name&sort_order=asc       |
| **Validation**   | ✅ Basic type checking      | ✅ Format validation (domain, IP)     |
| **Query Combo**  | ❌ Can't combine everything | ✅ All parameters work together       |
| **Response**     | Array of assets             | Paginated result with metadata        |
| **Code changes** | N/A                         | **Clean Architecture → Easy to add!** |

## New Architecture Components

```
session4-advanced-api/
├── internal/
│   ├── validator/                    # 🆕 NEW!
│   │   └── asset_validator.go        # Input validation logic
│   ├── storage/
│   │   ├── storage.go                # 🔄 Enhanced interface
│   │   └── postgres/
│   │       └── postgres.go           # 🔄 Dynamic query building
│   ├── service/
│   │   └── asset_service.go          # 🔄 Uses validator
│   └── handler/
│       └── asset_handler.go          # 🔄 Parse query params
├── cmd/server/main.go
├── go.mod
├── docker-compose.yml
└── README.md
```

---

## 🆕 Key Features

### 1. Pagination

Handle large datasets efficiently with page-based navigation.

**API:**

```bash
GET /assets?page=1&page_size=20
```

**Response:**

```json
{
  "data": [...],
  "total": 245,
  "page": 1,
  "page_size": 20,
  "total_pages": 13
}
```

**Defaults:**

- `page`: 1 (if not specified)
- `page_size`: 20 (if not specified, max 100)

### 2. Filtering

Filter by multiple criteria simultaneously.

**API:**

```bash
# Filter by type
GET /assets?type=domain

# Filter by status
GET /assets?status=active

# Combine filters
GET /assets?type=domain&status=active
```

**Supported Filters:**

- `type`: domain, ip, service
- `status`: active, inactive

### 3. Search

Partial name matching (case-insensitive).

**API:**

```bash
GET /assets?search=example
```

Matches: "example.com", "test-example.net", "my-example-service"

### 4. Sorting

Sort by any field in ascending or descending order.

**API:**

```bash
# Sort by name (A-Z)
GET /assets?sort_by=name&sort_order=asc

# Sort by creation date (newest first, default)
GET /assets?sort_by=created_at&sort_order=desc
```

**Sortable Fields:**

- `name`
- `type`
- `status`
- `created_at`
- `updated_at`

**Defaults:**

- `sort_by`: created_at
- `sort_order`: desc

### 5. Input Validation

Comprehensive validation with clear error messages.

**Domain Validation:**

```go
// ✅ Valid
example.com
subdomain.example.com
my-domain-123.net

// ❌ Invalid
invalid..com        // Double dot
-example.com        // Starts with hyphen
example..com        // Consecutive dots
```

**IP Validation:**

```go
// ✅ Valid IPv4
192.168.1.1
10.0.0.1

// ✅ Valid IPv6
2001:db8::1
fe80::1

// ❌ Invalid
999.999.999.999     // Out of range
not-an-ip           // Invalid format
```

**Service Validation:**

```go
// ✅ Valid
http://example.com
https://example.com:443
ssh
ftp

// ❌ Invalid
invalid service!!!  // Special characters
```

### 6. Combined Queries

All features work together!

```bash
curl "http://localhost:8080/assets?type=domain&status=active&search=example&page=2&page_size=10&sort_by=name&sort_order=asc"
```

This query:

- Filters: active domains containing "example"
- Pagination: page 2, 10 items per page
- Sorting: alphabetically by name

---

## 📋 Complete API Reference

### Health Check

```bash
GET /health
```

### Create Asset

```bash
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "type": "domain"
  }'
```

### List Assets (with all options)

```bash
curl "http://localhost:8080/assets?page=1&page_size=20&type=domain&status=active&search=example&sort_by=name&sort_order=asc"
```

**Query Parameters:**

- `page` (int, default: 1): Page number
- `page_size` (int, default: 20, max: 100): Items per page
- `type` (string): Filter by type (domain,ip, service)
- `status` (string): Filter by status (active, inactive)
- `search` (string): Search in name field
- `sort_by` (string): Sort field (name, type, status, created_at, updated_at)
- `sort_order` (string): Sort direction (asc, desc)

### Get Single Asset

```bash
curl http://localhost:8080/assets/{id}
```

### Update Asset

```bash
curl -X PUT http://localhost:8080/assets/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "updated.com",
    "status": "inactive"
  }'
```

### Delete Asset

```bash
curl -X DELETE http://localhost:8080/assets/{id}
```

---

## 🎓 Teaching Flow (3 hours)

### Part 1: Review & Motivation (15 min)

1. **Review Session 3**
   - Show current API working
   - List all assets → "What if we have 10,000 assets?"
   - Show limitation: no pagination, fixed sorting

2. **Real-world Problem**
   - Show Google search results: pages!
   - Show Amazon product listings: filters + sort
   - "Modern APIs need these features"

### Part 2: Validator Package (30 min)

1. **Why Separate Validation?**
   - Code organization
   - Reusability
   - Single responsibility

2. **Code Walkthrough**
   - `internal/validator/asset_validator.go`
   - Show domain validation with regex
   - Show IP validation with net.ParseIP
   - Show service validation

3. **Demo**
   - Test valid domain: ✅ accepted
   - Test invalid domain: ❌ rejected with clear error
   - Show error message quality

**Key Teaching Point:**

> "Validation at the boundary! Catch errors early with clear messages."

### Part 3: Storage Interface Enhancement (30 min)

1. **QueryParams Struct**
   - Show old vs new interface
   - Explain consolidation benefit

2. **Dynamic Query Building**
   - Open `internal/storage/postgres/postgres.go`
   - Explain WHERE 1=1 pattern
   - Show parameter counting ($1, $2, $3...)
   - Highlight security: prepared statements

3. **Pagination Implementation**
   - LIMIT and OFFSET
   - Calculate offset: (page - 1) \* pageSize
   - Show count query for total

**Key Teaching Point:**

> "Never concatenate user input into SQL! Always use prepared statements."

### Part 4: Service Layer Updates (20 min)

1. **Validator Integration**
   - Service now has validator field
   - Validation before business logic

2. **ListAssets Method**
   - Replaces GetAll, Filter, Search
   - Sets defaults
   - Validates parameters
   - Single method, many use cases

**Key Teaching Point:**

> "Service sets sensible defaults. Handler just passes data."

### Part 5: Handler Enhancements (25 min)

1. **Query Parameter Parsing**
   - parseIntParam helper
   - URL query parsing
   - Building QueryParams struct

2. **Error Handling**
   - mapErrorToStatus function
   - Validation errors → 400
   - Not found → 404
   - Server errors → 500

**Key Teaching Point:**

> "Handler translates between HTTP and domain. No business logic here!"

### Part 6: Live Demo (45 min)

**Setup (5 min):**

```bash
docker-compose up -d
go run cmd/server/main.go
```

**Create Test Data (10 min):**

```bash
# Create 20 different assets for testing
for i in {1..10}; do
  curl -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"example$i.com\",\"type\":\"domain\"}"
done

for i in {1..5}; do
  curl -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"192.168.1.$i\",\"type\":\"ip\"}"
done

for i in {1..5}; do
  curl -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"http://service$i.com\",\"type\":\"service\"}"
done
```

**Demo Scenarios (30 min):**

1. **Pagination** (5 min)

   ```bash
   # Page 1
   curl "http://localhost:8080/assets?page=1&page_size=5"
   # Show: 5 items, total, total_pages

   # Page 2
   curl "http://localhost:8080/assets?page=2&page_size=5"
   # Show: different 5 items
   ```

2. **Filtering** (5 min)

   ```bash
   # Only domains
   curl "http://localhost:8080/assets?type=domain"

   # Only IPs
   curl "http://localhost:8080/assets?type=ip"

   # Active domains
   curl "http://localhost:8080/assets?type=domain&status=active"
   ```

3. **Search** (5 min)

   ```bash
   # Search "example"
   curl "http://localhost:8080/assets?search=example"
   # Show: only example1.com, example2.com, etc.

   # Search "192"
   curl "http://localhost:8080/assets?search=192"
   # Show: only IPs
   ```

4. **Sorting** (5 min)

   ```bash
   # Alphabetical
   curl "http://localhost:8080/assets?sort_by=name&sort_order=asc"

   # Reverse alphabetical
   curl "http://localhost:8080/assets?sort_by=name&sort_order=desc"
   ```

5. **Combined** (5 min)

   ```bash
   # The ultimate query
   curl "http://localhost:8080/assets?type=domain&search=example&page=1&page_size=3&sort_by=name&sort_order=asc"
   ```

6. **Validation** (5 min)

   ```bash
   # Invalid domain
   curl -X POST http://localhost:8080/assets \
     -H "Content-Type: application/json" \
     -d '{"name":"invalid..com","type":"domain"}'
   # Show: 400 with clear error

   # Invalid IP
   curl -X POST http://localhost:8080/assets \
     -H "Content-Type: application/json" \
     -d '{"name":"999.999.999.999","type":"ip"}'
   # Show: 400 with clear error
   ```

### Part 7: Architecture Discussion (15 min)

**Comparison Table:**

| Layer     | Session 3        | Session 4                         |
| --------- | ---------------- | --------------------------------- |
| Model     | Same             | Same                              |
| Storage   | Simple methods   | QueryParams, dynamic queries      |
| Service   | Basic validation | Validator package, defaults       |
| Handler   | No query parsing | Parse & validate query params     |
| Validator | N/A              | **NEW!** Type-specific validation |

**Key Points:**

- Clean Architecture makes this easy!
- Each layer has clear responsibility
- Easy to test each layer independently
- No changes to model layer!

### Part 8: Q&A & Homework (10 min)

**Common Questions:**

- Q: "Why not validate in handler?"
  A: "Handler is HTTP concerns only. Validation is business logic."

- Q: "What if database is slow with complex queries?"
  A: "Add indexes! We already have them in migration."

- Q: "Can we add more filters?"
  A: "Yes! Just add to QueryParams and update query building."

**Homework:**

1. Add date range filtering (`created_after`, `created_before`)
2. Add count-only endpoint (`GET /assets/count`)
3. Add bulk create (`POST /assets/bulk`)
4. Write unit tests for validator
5. Write integration tests for queries

---

## 🔑 Key Concepts

### 1. Dynamic Query Building

```go
query := "SELECT * FROM assets WHERE 1=1"
args := []interface{}{}

if params.Type != "" {
    query += fmt.Sprintf(" AND type = $%d", len(args)+1)
    args = append(args, params.Type)
}
// Add more conditions...
```

**Why `WHERE 1=1`?**

- Simplifies dynamic AND conditions
- Don't need to track if first condition
- Clean code pattern

### 2. Pagination Math

```go
offset := (page - 1) * pageSize
totalPages := int(total) / pageSize
if int(total) % pageSize != 0 {
    totalPages++
}
```

**Examples:**

- Page 1, Size 10 → Offset 0, Limit 10
- Page 2, Size 10 → Offset 10, Limit 10
- Total 25, Size 10 → 3 pages

### 3. Whitelist Validation

```go
validFields := map[string]bool{
    "name": true,
    "created_at": true,
}

if !validFields[sortBy] {
    sortBy = "created_at" // Use default
}
```

**Critical for Security:**

- Prevents SQL injection
- User can't specify arbitrary columns
- Always validate sort/filter fields!

### 4. Error Mapping

```go
func mapErrorToStatus(err error) int {
    switch err {
    case model.ErrNotFound:
        return http.StatusNotFound
    case model.ErrInvalidInput:
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

**Separation of Concerns:**

- Domain errors (model package)
- HTTP status codes (handler)
- Service doesn't know about HTTP!

---

## 🐛 Common Issues & Solutions

### Issue 1: Query Returns No Results

**Symptom:** Empty array despite having data

**Check:**

```bash
# Verify filters
curl "http://localhost:8080/assets?type=domains" # Wrong! Should be "domain"

# Check search
curl "http://localhost:8080/assets?search=EXAMPLE" # Should still work (case-insensitive)
```

**Solution:** Verify filter values match constants (domain, not domains)

### Issue 2: Pagination Shows Wrong Total

**Symptom:** `total_pages` incorrect

**Debug:**

```sql
-- Check count query
SELECT COUNT(*) FROM assets WHERE type = 'domain';

-- Check data query
SELECT * FROM assets WHERE type = 'domain' LIMIT 10 OFFSET 0;
```

**Solution:** Ensure count query has same filters as data query

### Issue 3: Sort Field Rejected

**Symptom:** "invalid sort field" error

**Check:** Whitelist in validator and storage

```go
// Must be in both places!
validSortFields := map[string]bool{
    "name": true,
    "type": true,
    "created_at": true,
}
```

### Issue 4: Validation Too Strict

**Symptom:** Valid domains rejected

**Example:**

```bash
# This might fail if regex too strict
curl -X POST http://localhost:8080/assets \
  -d '{"name":"my-123-domain.com","type":"domain"}'
```

**Solution:** Review regex pattern, ensure RFC compliance

---

## 📚 Resources

### Further Reading

- [pagination best practices](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)
- [SQL injection prevention](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html)
- [Input validation strategies](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html)
- [Domain name validation](https://en.wikipedia.org/wiki/Domain_name#Domain_name_syntax)

### Tools

- [Postman](https://www.postman.com/) - API testing
- [httpie](https://httpie.io/) - User-friendly curl alternative
- [jq](https://stedolan.github.io/jq/) - JSON processor for terminal

---

## ✅ Success Criteria

Students should be able to:

- [ ] Explain why pagination is important
- [ ] Implement dynamic query building safely
- [ ] Write type-specific validation
- [ ] Combine multiple query parameters
- [ ] Map domain errors to HTTP status codes
- [ ] Explain SQL injection prevention
- [ ] Test API with complex queries

---

**Next Session:** Testing & Quality (Unit tests, integration tests, mocks)
"updated_at": true,
}

    if !validSortFields[sortBy] {
        sortBy = "created_at"
    }

    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }

    query := fmt.Sprintf("SELECT * FROM assets ORDER BY %s %s", sortBy, sortOrder)
    // ...

}

````

### 3. Advanced Validation

```go
// internal/validator/asset_validator.go
package validator

type AssetValidator struct{}

func (v *AssetValidator) ValidateCreate(req CreateAssetRequest) error {
    // Name validation
    if len(req.Name) == 0 {
        return errors.New("name is required")
    }
    if len(req.Name) > 255 {
        return errors.New("name too long (max 255 characters)")
    }

    // Type-specific validation
    switch req.Type {
    case "domain":
        if !isDomainValid(req.Name) {
            return errors.New("invalid domain format")
        }
    case "ip":
        if !isIPValid(req.Name) {
            return errors.New("invalid IP address")
        }
    case "service":
        if !isServiceValid(req.Name) {
            return errors.New("invalid service format")
        }
    }

    return nil
}

func isDomainValid(domain string) bool {
    // Must have at least one dot
    // No spaces, special characters
    pattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
    matched, _ := regexp.MatchString(pattern, domain)
    return matched
}

func isIPValid(ip string) bool {
    return net.ParseIP(ip) != nil
}
````

### 4. Complex Filtering

**Multiple filters combined:**

```bash
GET /assets?type=domain&status=active&search=example&page=1&page_size=10&sort_by=name
```

**Implementation:**

```go
type FilterParams struct {
    Type      string
    Status    string
    Search    string
    Page      int
    PageSize  int
    SortBy    string
    SortOrder string
}

func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
    params := FilterParams{
        Type:      r.URL.Query().Get("type"),
        Status:    r.URL.Query().Get("status"),
        Search:    r.URL.Query().Get("search"),
        Page:      getIntParam(r, "page", 1),
        PageSize:  getIntParam(r, "page_size", 20),
        SortBy:    r.URL.Query().Get("sort_by"),
        SortOrder: r.URL.Query().Get("sort_order"),
    }

    result, err := h.service.FilterAssets(params)
    // ...
}
```

### 5. Bulk Operations

```go
// Bulk create
POST /assets/bulk
{
  "assets": [
    {"name": "example1.com", "type": "domain"},
    {"name": "example2.com", "type": "domain"}
  ]
}

// Bulk delete
DELETE /assets/bulk
{
  "ids": ["uuid1", "uuid2", "uuid3"]
}
```

### 6. Error Response Enhancement

```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Code    string                 `json:"code"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// Example response
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "name": "name is required",
    "type": "invalid type"
  }
}
```

## Performance Considerations

### Database Indexes

```sql
-- Composite indexes for common queries
CREATE INDEX idx_assets_type_status ON assets(type, status);
CREATE INDEX idx_assets_name_search ON assets(name text_pattern_ops);

-- For pagination
CREATE INDEX idx_assets_created_at_id ON assets(created_at DESC, id);
```

### Query Optimization

```go
// ❌ BAD - N+1 problem
for _, asset := range assets {
    details := getAssetDetails(asset.ID) // Multiple queries
}

// ✅ GOOD - Single query with JOIN
assets := getAssetsWithDetails() // One query
```

## Testing Scenarios

```bash
# Test pagination
curl "http://localhost:8080/assets?page=1&page_size=5"
curl "http://localhost:8080/assets?page=2&page_size=5"

# Test sorting
curl "http://localhost:8080/assets?sort_by=name&sort_order=asc"
curl "http://localhost:8080/assets?sort_by=created_at&sort_order=desc"

# Test combined filters
curl "http://localhost:8080/assets?type=domain&status=active&search=example&page=1&sort_by=name"

# Test validation
curl -X POST http://localhost:8080/assets \
  -d '{"name":"", "type":"domain"}' # Should fail

curl -X POST http://localhost:8080/assets \
  -d '{"name":"invalid..domain", "type":"domain"}' # Should fail
```

## Homework

1. **Add date range filtering**

   ```
   GET /assets?created_after=2026-01-01&created_before=2026-12-31
   ```

2. **Add export functionality**

   ```
   GET /assets/export?format=csv
   GET /assets/export?format=json
   ```

3. **Add field selection**

   ```
   GET /assets?fields=id,name,type
   ```

4. **Add aggregation endpoint**
   ```
   GET /assets/stats
   {
     "total": 100,
     "by_type": {"domain": 60, "ip": 30, "service": 10},
     "by_status": {"active": 80, "inactive": 20}
   }
   ```

## Resources

- [SQL Performance Explained](https://use-the-index-luke.com/)
- [Go Validator Package](https://github.com/go-playground/validator)
- [REST API Best Practices](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
