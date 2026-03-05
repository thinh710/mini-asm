package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"mini-asm/internal/config"
	"mini-asm/internal/model"
	"mini-asm/internal/storage"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresStorage implements the Storage interface using PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func NewPostgresStorageFromConfig(config *config.PostgresConfig) (*PostgresStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresStorage{db: db}, nil
}

// Create inserts a new asset into the database
func (p *PostgresStorage) Create(asset *model.Asset) error {
	query := `
		INSERT INTO assets (id, name, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.db.Exec(
		query,
		asset.ID,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.CreatedAt,
		asset.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	return nil
}

// GetAll retrieves assets with filtering, sorting, and pagination
func (p *PostgresStorage) GetAll(params storage.QueryParams) (*storage.PaginatedResult, error) {
	// Build base query with filters
	query, args := p.buildQuery(params)

	// Add sorting
	query += p.buildOrderBy(params.SortBy, params.SortOrder)

	// Add pagination
	offset := (params.Page - 1) * params.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.PageSize, offset)

	// Execute query
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	// Scan results
	assets := []*model.Asset{}
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Get total count for pagination metadata
	total, err := p.Count(params)
	if err != nil {
		return nil, fmt.Errorf("failed to count assets: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	return &storage.PaginatedResult{
		Data:       assets,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Count returns the total number of assets matching the filters
func (p *PostgresStorage) Count(params storage.QueryParams) (int64, error) {
	// Build count query with same filters
	countQuery, args := p.buildCountQuery(params)

	var total int64
	err := p.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count assets: %w", err)
	}

	return total, nil
}

// buildQuery constructs the SELECT query with filters
func (p *PostgresStorage) buildQuery(params storage.QueryParams) (string, []interface{}) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Filter by type
	if params.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, params.Type)
		argCount++
	}

	// Filter by status
	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, params.Status)
		argCount++
	}

	// Search by name (case-insensitive partial match)
	if params.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	return query, args
}

// buildCountQuery constructs the COUNT query with filters
func (p *PostgresStorage) buildCountQuery(params storage.QueryParams) (string, []interface{}) {
	query := "SELECT COUNT(*) FROM assets WHERE 1=1"

	args := []interface{}{}
	argCount := 1

	// Apply same filters as main query
	if params.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, params.Type)
		argCount++
	}

	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, params.Status)
		argCount++
	}

	if params.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	return query, args
}

// buildOrderBy constructs the ORDER BY clause
func (p *PostgresStorage) buildOrderBy(sortBy, sortOrder string) string {
	// Default sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Whitelist validation is done in validator, but double-check here
	validFields := map[string]bool{
		"name":       true,
		"type":       true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}

	if !validFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}

// GetByID retrieves a single asset by its ID
func (p *PostgresStorage) GetByID(id string) (*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE id = $1
	`

	asset := &model.Asset{}
	err := p.db.QueryRow(query, id).Scan(
		&asset.ID,
		&asset.Name,
		&asset.Type,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return asset, nil
}

// Update modifies an existing asset in the database
func (p *PostgresStorage) Update(id string, asset *model.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, type = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := p.db.Exec(
		query,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Delete removes an asset from the database
func (p *PostgresStorage) Delete(id string) error {
	query := `DELETE FROM assets WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

/*
🎓 TEACHING NOTES:

=== SESSION 4 ENHANCEMENTS ===

1. Unified GetAll Method:
   - Replaces separate GetAll(), Filter(), Search() methods
   - Single method handles all query variations
   - More maintainable and flexible

2. Dynamic Query Building:
   - buildQuery() constructs WHERE clause dynamically
   - Only adds conditions for provided parameters
   - Maintains prepared statement security ($1, $2, ...)

3. Query Construction Pattern:
   ```
   WHERE 1=1              ← Always true, simplifies logic
   AND type = $1          ← Add condition if param present
   AND status = $2        ← Add condition if param present
   AND name ILIKE $3      ← Add condition if param present
   ORDER BY created_at DESC
   LIMIT $4 OFFSET $5
   ```

4. Pagination Implementation:
   - LIMIT: how many results to return
   - OFFSET: how many results to skip
   - Offset = (Page - 1) * PageSize
   - Example: Page 2, Size 10 → OFFSET 10

5. Sorting:
   - buildOrderBy() creates ORDER BY clause
   - Whitelist validation (critical for security!)
   - Default: ORDER BY created_at DESC
   - Never concatenate user input directly!

6. Count Query:
   - Separate query to get total count
   - Same filters as main query
   - Used for pagination metadata

7. Parameter Counting:
   - Track $1, $2, $3... positions
   - argCount increments for each parameter
   - Important: LIMIT and OFFSET use next available numbers

SECURITY NOTES:

1. SQL Injection Prevention:
   - ✅ Always use prepared statements ($1, $2)
   - ✅ Whitelist sort fields
   - ✅ Validate parameters before query
   - ❌ Never: fmt.Sprintf("ORDER BY %s", userInput)

2. Defense in Depth:
   - Validation in validator layer
   - Double-check in storage layer
   - Prepared statements in database

EXAMPLE QUERIES:

1. All assets, page 1:
   ```sql
   SELECT ... FROM assets WHERE 1=1
   ORDER BY created_at DESC
   LIMIT 20 OFFSET 0
   ```

2. Filter by type, page 2:
   ```sql
   SELECT ... FROM assets WHERE 1=1 AND type = $1
   ORDER BY created_at DESC
   LIMIT 20 OFFSET 20
   ```
   Args: ["domain", 20, 20]

3. Filter + Search + Sort:
   ```sql
   SELECT ... FROM assets WHERE 1=1
   AND type = $1 AND name ILIKE $2
   ORDER BY name ASC
   LIMIT 10 OFFSET 0
   ```
   Args: ["domain", "%example%", 10, 0]

PERFORMANCE CONSIDERATIONS:

1. Indexes Help:
   - idx_assets_type → speeds up type filter
   - idx_assets_name → speeds up ILIKE search
   - idx_assets_created_at → speeds up default sort

2. ILIKE Performance:
   - Case-insensitive search
   - Can be slow on large datasets
   - Consider full-text search for production

3. Count Query:
   - Can be expensive on large tables
   - Consider caching count results
   - Or estimate with EXPLAIN

COMPARISON WITH SESSION 3:

Session 3:
  - Separate methods for each operation
  - No pagination
  - No sorting options
  - Harder to combine filters

Session 4:
  - Single unified method
  - Built-in pagination
  - Flexible sorting
  - Easy to combine any filters

Example improvement:
  Session 3: Can't get "active domains" with pagination
  Session 4: params := QueryParams{Type: "domain", Status: "active", Page: 1}

TESTING STRATEGIES:

1. Unit tests with mock DB
2. Integration tests with real DB
3. Test edge cases:
   - Empty results
   - Large datasets
   - Special characters in search
   - Invalid sort fields (should use default)

DEMO POINTS:

1. Show query building with fmt.Printf
2. Show parameter array construction
3. Test with different filter combinations
4. Show pagination working (page 1, 2, 3)
5. Show sorting (asc vs desc)
6. Show combined: filter + search + sort + page
*/
