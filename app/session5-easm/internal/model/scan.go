package model

import "time"

// ScanType represents the type of scan being performed
type ScanType string

const (
	ScanTypeSubdomain ScanType = "subdomain"
	ScanTypeDNS       ScanType = "dns"
	ScanTypeWHOIS     ScanType = "whois"
	ScanTypePort      ScanType = "port"
	ScanTypeASN       ScanType = "asn"
	ScanTypeSSL       ScanType = "ssl"
)

// ScanStatus represents the status of a scan
type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
	ScanStatusPartial   ScanStatus = "partial"
)

// ScanJob represents a scan task
type ScanJob struct {
	ID        string     `json:"id"`         // UUID
	AssetID   string     `json:"asset_id"`   // Foreign key to assets table
	ScanType  ScanType   `json:"scan_type"`  // Type of scan
	Status    ScanStatus `json:"status"`     // Current status
	StartedAt time.Time  `json:"started_at"` // When scan started
	EndedAt   *time.Time `json:"ended_at"`   // When scan completed (nullable)
	Error     string     `json:"error"`      // Error message if failed
	Results   int        `json:"results"`    // Number of results found
	CreatedAt time.Time  `json:"created_at"`
}

// Subdomain represents a discovered subdomain
type Subdomain struct {
	ID        string    `json:"id"`          // UUID
	AssetID   string    `json:"asset_id"`    // Parent domain asset
	ScanJobID string    `json:"scan_job_id"` // Which scan discovered this
	Name      string    `json:"name"`        // e.g., "api.example.com"
	Source    string    `json:"source"`      // How it was discovered
	IsActive  bool      `json:"is_active"`   // Reachable or not
	CreatedAt time.Time `json:"created_at"`
}

// DNSRecord represents a DNS record
type DNSRecord struct {
	ID         string    `json:"id"`          // UUID
	AssetID    string    `json:"asset_id"`    // Domain or subdomain asset
	ScanJobID  string    `json:"scan_job_id"` // Which scan discovered this
	RecordType string    `json:"record_type"` // A, AAAA, CNAME, MX, NS, TXT, SOA
	Name       string    `json:"name"`        // Record name
	Value      string    `json:"value"`       // Record value
	TTL        int       `json:"ttl"`         // Time to live
	CreatedAt  time.Time `json:"created_at"`
}

// WHOISRecord represents WHOIS information
type WHOISRecord struct {
	ID          string     `json:"id"`           // UUID
	AssetID     string     `json:"asset_id"`     // Domain asset
	ScanJobID   string     `json:"scan_job_id"`  // Which scan discovered this
	Registrar   string     `json:"registrar"`    // Domain registrar
	CreatedDate *time.Time `json:"created_date"` // Domain creation date (nullable)
	ExpiryDate  *time.Time `json:"expiry_date"`  // Domain expiry date (nullable)
	NameServers string     `json:"name_servers"` // JSON array of nameservers
	Status      string     `json:"status"`       // Domain status
	Emails      string     `json:"emails"`       // JSON array of contact emails
	RawData     string     `json:"raw_data"`     // Full WHOIS response
	CreatedAt   time.Time  `json:"created_at"`
}

// IsValidScanType checks if the given scan type is valid
func IsValidScanType(t ScanType) bool {
	switch t {
	case ScanTypeSubdomain, ScanTypeDNS, ScanTypeWHOIS, ScanTypePort, ScanTypeASN, ScanTypeSSL:
		return true
	}
	return false
}

// IsValidScanStatus checks if the given scan status is valid
func IsValidScanStatus(s ScanStatus) bool {
	switch s {
	case ScanStatusPending, ScanStatusRunning, ScanStatusCompleted, ScanStatusFailed, ScanStatusPartial:
		return true
	}
	return false
}

/*
🎓 TEACHING NOTES - Session 5: EASM Scanning

=== WHAT IS EASM (External Attack Surface Management)? ===

EASM helps organizations discover and monitor their external-facing assets:
- Domains and subdomains
- IP addresses
- Open ports and services
- SSL certificates
- DNS records
- WHOIS information

Why important?
- Security teams need to know: "What can attackers see?"
- Discovery of shadow IT and forgotten assets
- Continuous monitoring for new exposures

=== DOMAIN MODEL DESIGN ===

1. ScanJob (Job Tracking):
   Purpose: Track each scan operation
   - Async operations: scans take time
   - Status tracking: pending → running → completed/failed
   - Results count: quick summary without querying results
   - Error handling: store error messages for debugging

   Design Pattern: Job/Task Pattern
   - Used in background processing systems
   - Allows queuing and scheduling
   - Can monitor progress

2. Subdomain (Discovery Results):
   Purpose: Store discovered subdomains
   - Name: Full subdomain (e.g., "api.example.com")
   - Source: How discovered (DNS bruteforce, certificate transparency, web scraping)
   - IsActive: Whether currently reachable (important for asset inventory)
   - Foreign keys: AssetID (parent domain), ScanJobID (provenance tracking)

3. DNSRecord (DNS Information):
   Purpose: Store DNS records for domains
   - RecordType: A, AAAA, CNAME, MX, NS, TXT, SOA, etc.
   - Multiple records per domain (common)
   - TTL: Cache duration (useful for monitoring changes)

   Real-world use:
   - Find mail servers (MX records)
   - Discover CDN usage (CNAME)
   - Find IP addresses (A/AAAA records)

4. WHOISRecord (Registration Info):
   Purpose: Domain registration details
   - Registrar: Who manages the domain
   - Dates: Track expiration (security risk if domain expires!)
   - NameServers: DNS infrastructure
   - Emails: Contact information (potential phishing targets)
   - RawData: Keep original response for analysis

=== DATABASE DESIGN CONSIDERATIONS ===

1. Relationships:
   Asset (1) → (N) ScanJob
   Asset (1) → (N) Subdomain
   Asset (1) → (N) DNSRecord
   Asset (1) → (N) WHOISRecord
   ScanJob (1) → (N) Results

2. Normalization vs Denormalization:
   - Normalized: Store nameservers in separate table
   - Denormalized: Store as JSON array in string
   - We choose denormalized for simplicity (teaching project)
   - Production: might normalize for querying

3. Nullable Fields:
   - EndedAt: null while scan running
   - ExpiryDate: might not be parseable from WHOIS
   - Using pointers (*time.Time) for nullable timestamps

4. Indexes (not shown in model, added in migration):
   - asset_id: frequently filtered
   - scan_job_id: joining results to jobs
   - scan_type + status: dashboard queries

=== SCANNING FLOW (HIGH-LEVEL) ===

Client Request:
  POST /assets/{id}/scan
  Body: {"scan_type": "subdomain"}

Server Flow:
  1. Validate asset exists
  2. Create ScanJob (status: pending)
  3. Return job ID immediately (async pattern)
  4. Start scan in background
  5. Update job status: running
  6. Perform scan (calls external tools/APIs)
  7. Store results in respective tables
  8. Update job: completed + result count
  9. Handle errors: update job status to failed

Client Polling:
  GET /scan-jobs/{id}
  Response: {"status": "running", "results": 5, ...}

Why Async?
- Scans can take minutes
- Don't block HTTP request (timeout)
- Better UX (progress updates)
- Can parallelize multiple scans

=== COMPARISON WITH PREVIOUS SESSIONS ===

Session 2-4: CRUD on single entity (Asset)
Session 5:
  - Multiple related entities
  - Async operations (scan jobs)
  - External API calls (DNS, WHOIS)
  - More complex business logic

New concepts:
  - Foreign keys and relationships
  - Job/task pattern
  - Background processing
  - External integrations

=== SECURITY CONSIDERATIONS ===

1. Input Validation:
   - Validate domain format (prevent injection)
   - Rate limiting (prevent abuse)
   - Sanitize inputs before shell commands

2. External Calls:
   - Timeout handling (unresponsive services)
   - Error handling (service down)
   - Respect rate limits of external APIs

3. Data Privacy:
   - WHOIS may contain PII (emails, names)
   - Consider GDPR implications
   - Secure storage and access control

=== IMPLEMENTATION NOTES ===

For teaching purposes, we implement:
- ✅ WHOIS: Simple TCP/HTTP calls
- ✅ DNS: Using net package (standard library)
- ✅ Subdomain: Basic wordlist bruteforce

Left as exercises:
- Port scanning (nmap integration)
- ASN lookup (whois servers)
- SSL certificate analysis (TLS handshake)
- Advanced subdomain (certificate transparency, web scraping)

This teaches:
- How to structure scan systems
- External API integration patterns
- Async job processing
- Data modeling for security tools

Students will complete the remaining scan types using the same patterns!
*/
