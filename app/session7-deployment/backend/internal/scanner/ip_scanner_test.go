package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

// TestIPScanner_Type kiểm tra scan type trả về đúng
func TestIPScanner_Type(t *testing.T) {
	s := NewIPScanner()
	if s.Type() != model.ScanTypeIP {
		t.Errorf("expected ScanTypeIP, got %s", s.Type())
	}
}

// TestIPScanner_RejectsNonIPAsset kiểm tra từ chối asset không phải IP
func TestIPScanner_RejectsNonIPAsset(t *testing.T) {
	s := NewIPScanner()
	asset := &model.Asset{
		Name: "example.com",
		Type: model.TypeDomain,
	}
	_, err := s.Scan(asset)
	if err == nil {
		t.Error("expected error for domain asset, got nil")
	}
}

// TestIPScanner_RejectsServiceAsset kiểm tra từ chối service asset
func TestIPScanner_RejectsServiceAsset(t *testing.T) {
	s := NewIPScanner()
	asset := &model.Asset{
		Name: "https://example.com",
		Type: model.TypeService,
	}
	_, err := s.Scan(asset)
	if err == nil {
		t.Error("expected error for service asset, got nil")
	}
}

// TestIPScanner_AcceptsIPAsset kiểm tra chấp nhận IP asset hợp lệ
func TestIPScanner_AcceptsIPAsset(t *testing.T) {
	s := NewIPScanner()
	asset := &model.Asset{
		Name: "8.8.8.8",
		Type: model.TypeIP,
	}
	result, err := s.Scan(asset)
	if err != nil {
		t.Skipf("skipping: network unavailable (%v)", err)
	}
	if result == nil {
		t.Error("expected result, got nil")
	}
	if result.IPAddress != "8.8.8.8" {
		t.Errorf("expected ip 8.8.8.8, got %s", result.IPAddress)
	}
}

// TestIPScanner_ResultFields kiểm tra kết quả có đủ fields
func TestIPScanner_ResultFields(t *testing.T) {
	s := NewIPScanner()
	asset := &model.Asset{
		Name: "1.1.1.1",
		Type: model.TypeIP,
	}
	result, err := s.Scan(asset)
	if err != nil {
		t.Skipf("skipping: network unavailable (%v)", err)
	}
	if result.Country == "" {
		t.Error("expected Country to be non-empty")
	}
	if result.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

// BenchmarkIPScanner_Scan đo hiệu năng IP scan
func BenchmarkIPScanner_Scan(b *testing.B) {
	s := NewIPScanner()
	asset := &model.Asset{
		Name: "8.8.8.8",
		Type: model.TypeIP,
	}
	for i := 0; i < b.N; i++ {
		s.Scan(asset)
	}
}
