package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

// TestPortScanner_Type kiểm tra scan type trả về đúng
func TestPortScanner_Type(t *testing.T) {
	s := NewPortScanner()
	if s.Type() != model.ScanTypePort {
		t.Errorf("expected ScanTypePort, got %s", s.Type())
	}
}

// TestPortScanner_BlocksPublicIP kiểm tra public IP bị block
func TestPortScanner_BlocksPublicIP(t *testing.T) {
	s := NewPortScanner()
	asset := &model.Asset{
		Name: "8.8.8.8",
		Type: model.TypeIP,
	}
	_, err := s.Scan(asset)
	if err == nil {
		t.Error("expected authorization error for public IP")
	}
}

// TestPortScanner_BlocksAnotherPublicIP kiểm tra thêm IP công cộng khác
func TestPortScanner_BlocksAnotherPublicIP(t *testing.T) {
	s := NewPortScanner()
	asset := &model.Asset{
		Name: "1.1.1.1",
		Type: model.TypeIP,
	}
	_, err := s.Scan(asset)
	if err == nil {
		t.Error("expected authorization error for public IP 1.1.1.1")
	}
}

// TestPortScanner_IsAuthorized_Localhost kiểm tra localhost được phép
func TestPortScanner_IsAuthorized_Localhost(t *testing.T) {
	s := NewPortScanner()
	tests := []struct {
		ip   string
		want bool
	}{
		{"127.0.0.1", true},
		{"localhost", true},
		{"::1", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"192.0.2.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			got := s.isAuthorized(tt.ip)
			if got != tt.want {
				t.Errorf("isAuthorized(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

// TestPortScanner_AllowsLocalhost kiểm tra scan localhost thành công
func TestPortScanner_AllowsLocalhost(t *testing.T) {
	s := NewPortScanner()
	asset := &model.Asset{
		Name: "127.0.0.1",
		Type: model.TypeIP,
	}
	results, err := s.Scan(asset)
	if err != nil {
		t.Errorf("expected no error for localhost, got: %v", err)
	}
	// results có thể rỗng nếu không có port nào mở, không sao
	t.Logf("Found %d open ports on localhost", len(results))
}

// TestPortScanner_RejectsNonIPAsset kiểm tra từ chối asset không hợp lệ
func TestPortScanner_RejectsNonIPAsset(t *testing.T) {
	s := NewPortScanner()
	asset := &model.Asset{
		Name: "example.com",
		Type: model.TypeDomain,
	}
	_, err := s.Scan(asset)
	if err == nil {
		t.Error("expected error for domain asset without localhost resolution")
	}
}

// TestPortScanner_ScanPort kiểm tra hàm scanPort
func TestPortScanner_ScanPort(t *testing.T) {
	s := NewPortScanner()

	// Port 80 trên localhost thường đóng trong môi trường dev
	// Chỉ kiểm tra hàm không panic
	result := s.scanPort("127.0.0.1", 99999) // port không hợp lệ
	if result {
		t.Error("port 99999 should not be open")
	}
}

// BenchmarkPortScanner_ScanPort đo hiệu năng scan 1 port
func BenchmarkPortScanner_ScanPort(b *testing.B) {
	s := NewPortScanner()
	for i := 0; i < b.N; i++ {
		s.scanPort("127.0.0.1", 80)
	}
}
