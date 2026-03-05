package handler

import (
	"encoding/json"
	"net/http"

	"mini-asm/internal/model"
	"mini-asm/internal/service"
)

// ScanHandler handles scan-related HTTP requests
type ScanHandler struct {
	scanService *service.ScanService
}

// NewScanHandler creates a new scan handler
func NewScanHandler(scanService *service.ScanService) *ScanHandler {
	return &ScanHandler{
		scanService: scanService,
	}
}

// StartScan initiates a scan for an asset
// POST /assets/{id}/scan
func (h *ScanHandler) StartScan(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from path
	assetID := r.PathValue("id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		ScanType model.ScanType `json:"scan_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Start scan
	job, err := h.scanService.StartScan(assetID, req.ScanType)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return job info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted (async operation)
	json.NewEncoder(w).Encode(job)
}

// GetScanJob retrieves scan job status
// GET /scan-jobs/{id}
func (h *ScanHandler) GetScanJob(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from path
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "job ID required", http.StatusBadRequest)
		return
	}

	// Get job
	job, err := h.scanService.GetScanJob(jobID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return job
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// GetScanResults retrieves results for a scan job
// GET /scan-jobs/{id}/results
func (h *ScanHandler) GetScanResults(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from path
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "job ID required", http.StatusBadRequest)
		return
	}

	// Get results
	results, err := h.scanService.GetScanResults(jobID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ListScanJobs retrieves all scan jobs for an asset
// GET /assets/{id}/scans
func (h *ScanHandler) ListScanJobs(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from path
	assetID := r.PathValue("id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	// Get jobs
	jobs, err := h.scanService.ListScanJobs(assetID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return jobs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// GetAssetSubdomains retrieves all subdomains for an asset
// GET /assets/{id}/subdomains
func (h *ScanHandler) GetAssetSubdomains(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from path
	assetID := r.PathValue("id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	// Get subdomains
	subdomains, err := h.scanService.GetAssetSubdomains(assetID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return subdomains
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subdomains)
}

// GetAssetDNS retrieves all DNS records for an asset
// GET /assets/{id}/dns
func (h *ScanHandler) GetAssetDNS(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from path
	assetID := r.PathValue("id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	// Get DNS records
	records, err := h.scanService.GetAssetDNSRecords(assetID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return records
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// GetAssetWHOIS retrieves WHOIS information for an asset
// GET /assets/{id}/whois
func (h *ScanHandler) GetAssetWHOIS(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from path
	assetID := r.PathValue("id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	// Get WHOIS record
	record, err := h.scanService.GetAssetWHOIS(assetID)
	if err != nil {
		status := mapErrorToStatus(err)
		http.Error(w, err.Error(), status)
		return
	}

	// Return record
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

/*
🎓 TEACHING NOTES - Scan Handler

=== ASYNC API PATTERN ===

HTTP 202 Accepted:
```
Client:
  POST /assets/{id}/scan
  Body: {"scan_type": "dns"}

Server:
  HTTP/1.1 202 Accepted
  Body: {
    "id": "job-123",
    "status": "pending",
    ...
  }
```

202 Accepted means:
- Request accepted for processing
- Processing not complete yet
- Client should poll for status

Alternative status codes:
- 200 OK: Synchronous (completed immediately)
- 201 Created: Resource created (not appropriate here)
- 202 Accepted: Async operation started ✅

=== POLLING PATTERN ===

Client flow:
```javascript
// 1. Start scan
const response = await fetch('/assets/123/scan', {
  method: 'POST',
  body: JSON.stringify({scan_type: 'dns'})
});
const job = await response.json();

// 2. Poll for completion
const jobId = job.id;
let status = 'pending';

while (status !== 'completed' && status !== 'failed') {
  await sleep(2000);  // Wait 2 seconds

  const statusResponse = await fetch(`/scan-jobs/${jobId}`);
  const jobStatus = await statusResponse.json();
  status = jobStatus.status;

  console.log(`Status: ${status}, Results: ${jobStatus.results}`);
}

// 3. Get results
if (status === 'completed') {
  const resultsResponse = await fetch(`/scan-jobs/${jobId}/results`);
  const results = await resultsResponse.json();
  console.log('Results:', results);
}
```

Improvements (production):
- Exponential backoff (poll every 2s, then 4s, then 8s)
- Max retries (don't poll forever)
- WebSocket updates (push instead of pull)
- Server-Sent Events (SSE)

=== PATH PARAMETERS ===

Go 1.22+ pattern:
```go
// Route: GET /assets/{id}/scan
assetID := r.PathValue("id")
```

Benefits:
- Clean URLs
- RESTful design
- Type-safe extraction

Old way (pre-1.22):
```go
// Had to use third-party router (gorilla/mux, chi)
// or manual path parsing
vars := mux.Vars(r)
assetID := vars["id"]
```

=== ERROR MAPPING ===

```go
status := mapErrorToStatus(err)
```

Consistent error handling:
- model.ErrNotFound → 404 Not Found
- Validation errors → 400 Bad Request
- Other errors → 500 Internal Server Error

From asset_handler.go:
```go
func mapErrorToStatus(err error) int {
    if errors.Is(err, model.ErrNotFound) {
        return http.StatusNotFound
    }
    // ... validation errors
    return http.StatusInternalServerError
}
```

=== RESTFUL ENDPOINT DESIGN ===

Resource-oriented URLs:

Assets (from Session 4):
```
POST   /assets              - Create asset
GET    /assets              - List assets
GET    /assets/{id}         - Get asset
PUT    /assets/{id}         - Update asset
DELETE /assets/{id}         - Delete asset
```

Scans (Session 5):
```
POST   /assets/{id}/scan           - Start scan (nested under asset)
GET    /assets/{id}/scans          - List scans for asset
GET    /scan-jobs/{id}             - Get scan job status
GET    /scan-jobs/{id}/results     - Get scan results

GET    /assets/{id}/subdomains     - Get subdomains for asset
GET    /assets/{id}/dns            - Get DNS records for asset
GET    /assets/{id}/whois          - Get WHOIS for asset
```

Pattern:
- Collection: /resources
- Item: /resources/{id}
- Nested: /resources/{id}/sub-resources
- Actions: POST /resources/{id}/action

=== JSON REQUEST/RESPONSE ===

Request parsing:
```go
var req struct {
    ScanType model.ScanType `json:"scan_type"`
}

json.NewDecoder(r.Body).Decode(&req)
```

Anonymous struct:
- No need for separate type
- Clear what fields are expected
- Used only in this handler

Response encoding:
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```

Automatic JSON marshalling:
- Struct → JSON
- Uses struct tags (`json:"field_name"`)

=== MULTIPLE RETRIEVAL ENDPOINTS ===

Why separate endpoints?

1. `/scan-jobs/{id}/results` - Get results for specific scan
   - Use case: "Show me what this scan found"
   - Returns: Subdomains OR DNS OR WHOIS (depends on scan type)

2. `/assets/{id}/subdomains` - Get all subdomains for asset
   - Use case: "Show me all subdomains ever found"
   - Returns: Subdomains from all scans, historical data

3. `/assets/{id}/dns` - Get all DNS records for asset
   - Use case: "Show me current DNS configuration"
   - Returns: Latest DNS scan results

Different perspectives on same data!

=== CONTENT NEGOTIATION ===

Setting headers:
```go
w.Header().Set("Content-Type", "application/json")
```

Tells client:
- Response is JSON
- Client can parse accordingly

Could support multiple formats:
```go
accept := r.Header.Get("Accept")
if accept == "text/csv" {
    // Return CSV
} else {
    // Return JSON
}
```

For simplicity: JSON only

=== STATUS CODE SEMANTICS ===

```
200 OK           - GET requests, success
201 Created      - POST creation success
202 Accepted     - Async operation started ✅
204 No Content   - DELETE success
400 Bad Request  - Invalid input
404 Not Found    - Resource doesn't exist
500 Server Error - Unexpected error
```

202 Accepted use cases:
- Scan operations (our case)
- Long-running processes
- Batch operations
- Email sending

=== COMPARISON WITH SESSION 4 ===

Session 4: Synchronous CRUD
- Request → Process → Response (immediate)

Session 5: Async operations
- Request → Queue → Response with job ID
- Client polls for status
- More complex but necessary for long operations

New HTTP concepts:
- 202 Accepted status
- Job tracking endpoints
- Polling pattern
- Multiple result endpoints

=== TESTING SCAN ENDPOINTS ===

cURL examples:

1. Start DNS scan:
```bash
curl -X POST http://localhost:8080/assets/123/scan \
  -H "Content-Type: application/json" \
  -d '{"scan_type": "dns"}'

# Response:
{
  "id": "job-456",
  "asset_id": "123",
  "scan_type": "dns",
  "status": "pending",
  "started_at": "2024-03-05T10:00:00Z",
  ...
}
```

2. Check status:
```bash
curl http://localhost:8080/scan-jobs/job-456

# Response:
{
  "id": "job-456",
  "status": "running",
  "results": 5,
  ...
}
```

3. Get results (once completed):
```bash
curl http://localhost:8080/scan-jobs/job-456/results

# Response:
[
  {
    "id": "dns-1",
    "record_type": "A",
    "name": "example.com",
    "value": "93.184.216.34",
    ...
  },
  ...
]
```

4. List all scans for asset:
```bash
curl http://localhost:8080/assets/123/scans

# Response:
[
  {"id": "job-456", "scan_type": "dns", "status": "completed", ...},
  {"id": "job-789", "scan_type": "whois", "status": "completed", ...}
]
```

5. Get all subdomains:
```bash
curl http://localhost:8080/assets/123/subdomains

# Response:
[
  {"name": "www.example.com", "is_active": true, ...},
  {"name": "api.example.com", "is_active": true, ...}
]
```

=== PRODUCTION IMPROVEMENTS ===

1. **Authentication**:
   ```go
   // Check API key
   apiKey := r.Header.Get("X-API-Key")
   if !isValidKey(apiKey) {
       http.Error(w, "unauthorized", 401)
       return
   }
   ```

2. **Rate Limiting**:
   ```go
   if !rateLimiter.Allow(clientIP) {
       http.Error(w, "too many requests", 429)
       return
   }
   ```

3. **Request ID**:
   ```go
   requestID := uuid.New().String()
   w.Header().Set("X-Request-ID", requestID)
   // Log all operations with this ID
   ```

4. **CORS Headers**:
   ```go
   w.Header().Set("Access-Control-Allow-Origin", "*")
   // Allow frontend to call API
   ```

5. **Webhook Notifications**:
   - When scan completes
   - POST to configured URL
   - No need for polling

=== KEY TAKEAWAYS ===

1. HTTP 202 for async operations
2. Polling pattern for status updates
3. RESTful nested resources
4. Multiple endpoints for different views
5. Consistent error handling
6. Path parameters for clean URLs
7. JSON as API format

Handler layer: HTTP boundary, translates requests to service calls!
*/
