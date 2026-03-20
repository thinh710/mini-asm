package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mini-asm/internal/model"
)

type IPScanner struct {
	client *http.Client
}

type IPScanResult struct {
	IPAddress   string    `json:"ip_address"`
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	City        string    `json:"city"`
	Region      string    `json:"region"`
	ISP         string    `json:"isp"`
	Org         string    `json:"org"`
	ASN         string    `json:"asn"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewIPScanner() *IPScanner {
	return &IPScanner{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *IPScanner) Type() model.ScanType {
	return model.ScanTypeIP
}

func (s *IPScanner) Scan(asset *model.Asset) (*IPScanResult, error) {
	if asset.Type != model.TypeIP {
		return nil, fmt.Errorf("ip scan requires IP asset, got: %s", asset.Type)
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", asset.Name)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ip lookup failed: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	getString := func(key string) string {
		if v, ok := data[key]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}

	result := &IPScanResult{
		IPAddress:   asset.Name,
		Country:     getString("country"),
		CountryCode: getString("countryCode"),
		City:        getString("city"),
		Region:      getString("regionName"),
		ISP:         getString("isp"),
		Org:         getString("org"),
		ASN:         getString("as"),
		CreatedAt:   time.Now(),
	}
	if lat, ok := data["lat"].(float64); ok {
		result.Latitude = lat
	}
	if lon, ok := data["lon"].(float64); ok {
		result.Longitude = lon
	}
	return result, nil
}
