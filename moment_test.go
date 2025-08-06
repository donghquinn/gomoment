package gomoment_test

import (
	"strings"
	"testing"
	"time"

	"github.com/donghquinn/gomoment"
)

// Test data for consistent testing
var (
	testTime      = time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)
	testTimeLocal = time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	testTimeWithMS = time.Date(2023, 5, 15, 14, 30, 45, 123000000, time.UTC)
)

func TestNewMoment_NoArgs(t *testing.T) {
	moment, err := gomoment.NewMoment()
	if err != nil {
		t.Fatalf("NewMoment() failed: %v", err)
	}
	
	// Should be close to current time (within 1 second)
	now := time.Now()
	diff := now.Sub(moment.Time())
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("NewMoment() time too different from now: got %v, want within 1s of %v", moment.Time(), now)
	}
}

func TestNewMoment_WithTime(t *testing.T) {
	moment, err := gomoment.NewMoment(testTime)
	if err != nil {
		t.Fatalf("NewMoment(time.Time) failed: %v", err)
	}
	
	if !moment.Time().Equal(testTime) {
		t.Errorf("NewMoment(time.Time) = %v, want %v", moment.Time(), testTime)
	}
}

func TestNewMoment_WithString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "ISO date",
			input:    "2023-05-15",
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "ISO datetime",
			input:    "2023-05-15 14:30:45",
			wantTime: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
		},
		{
			name:     "RFC3339",
			input:    "2023-05-15T14:30:45Z",
			wantTime: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
		},
		{
			name:     "RFC3339 with timezone",
			input:    "2023-05-15T14:30:45+09:00",
			wantTime: time.Date(2023, 5, 15, 5, 30, 45, 0, time.UTC), // Converted to UTC
		},
		{
			name:     "US format",
			input:    "05/15/2023",
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "US format with time",
			input:    "05/15/2023 14:30:45",
			wantTime: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
		},
		{
			name:    "invalid format",
			input:   "invalid-date-string",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moment, err := gomoment.NewMoment(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMoment(%q) expected error, got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("NewMoment(%q) failed: %v", tt.input, err)
			}
			
			if !moment.Time().Equal(tt.wantTime) {
				t.Errorf("NewMoment(%q) = %v, want %v", tt.input, moment.Time(), tt.wantTime)
			}
		})
	}
}

func TestNewMoment_InvalidArgType(t *testing.T) {
	_, err := gomoment.NewMoment(123)
	if err == nil {
		t.Error("NewMoment(int) expected error, got nil")
	}
	
	expectedMsg := "unsupported argument type"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("NewMoment(int) error = %v, want to contain %q", err, expectedMsg)
	}
}

func TestNow(t *testing.T) {
	moment := gomoment.Now()
	
	// Should be close to current time (within 1 second)
	now := time.Now()
	diff := now.Sub(moment.Time())
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Now() time too different from current time: got %v, want within 1s of %v", moment.Time(), now)
	}
}

func TestFormat_BasicTokens(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"4-digit year", "YYYY", "2023"},
		{"2-digit year", "YY", "23"},
		{"Month zero-padded", "MM", "05"},
		{"Month", "M", "5"},
		{"Day zero-padded", "DD", "15"},
		{"Day", "D", "15"},
		{"Hour 24-format zero-padded", "HH", "14"},
		{"Hour 24-format", "H", "14"},
		{"Hour 12-format zero-padded", "hh", "02"},
		{"Hour 12-format", "h", "2"},
		{"Minutes zero-padded", "mm", "30"},
		{"Minutes", "m", "30"},
		{"Seconds zero-padded", "ss", "45"},
		{"Seconds", "s", "45"},
		{"Milliseconds", "SSS", "000"},
		{"Timezone offset", "ZZ", "+0000"}, // UTC
		{"Timezone offset with colon", "Z", "+00:00"}, // UTC
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := moment.Format(tt.format)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.format, err)
			}
			if result != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.format, result, tt.want)
			}
		})
	}
}

func TestFormat_CommonFormats(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"ISO date", "YYYY-MM-DD", "2023-05-15"},
		{"ISO datetime", "YYYY-MM-DD HH:mm:ss", "2023-05-15 14:30:45"},
		{"US format", "MM/DD/YYYY", "05/15/2023"},
		{"US format short", "M/D/YY", "5/15/23"},
		{"12-hour format", "h:mm", "2:30"},
		{"24-hour format", "HH:mm", "14:30"},
		{"Full format with timezone", "YYYY-MM-DD HH:mm:ss Z", "2023-05-15 14:30:45 +00:00"},
		{"Custom format", "D MM YYYY h:mm", "15 05 2023 2:30"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := moment.Format(tt.format)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.format, err)
			}
			if result != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.format, result, tt.want)
			}
		})
	}
}

func TestFormat_InvalidTokens(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	tests := []struct {
		name   string
		format string
	}{
		{"Single invalid token", "YYYY-QQ-DD"},
		{"Multiple invalid tokens", "INVALID-FORMAT"},
		{"Mixed valid and invalid", "YYYY-XX-MM-YY"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := moment.Format(tt.format)
			if err == nil {
				t.Errorf("Format(%q) expected error, got nil", tt.format)
			}
			
			expectedMsg := "unsupported format tokens"
			if !strings.Contains(err.Error(), expectedMsg) {
				t.Errorf("Format(%q) error = %v, want to contain %q", tt.format, err, expectedMsg)
			}
		})
	}
}

func TestFormat_EdgeCases(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"Empty format", "", ""},
		{"Only separators", "---", "---"},
		{"Repeated tokens", "YYYY-YYYY", "2023-2023"},
		{"Mixed separators", "YYYY/MM-DD", "2023/05-15"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := moment.Format(tt.format)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.format, err)
			}
			if result != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.format, result, tt.want)
			}
		})
	}
}

func TestMust_Success(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	result := moment.Must("YYYY-MM-DD")
	expected := "2023-05-15"
	
	if result != expected {
		t.Errorf("Must(%q) = %q, want %q", "YYYY-MM-DD", result, expected)
	}
}

func TestMust_Panic(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must with invalid format should panic")
		}
	}()
	
	moment.Must("INVALID-FORMAT")
}

func TestParse_VariousFormats(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "Date only",
			input:    "2023-05-15",
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Slash format",
			input:    "2023/05/15",
			wantTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "With time",
			input:    "2023-05-15 14:30:45",
			wantTime: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
		},
		{
			name:     "Time only",
			input:    "14:30:45",
			wantTime: time.Date(0, 1, 1, 14, 30, 45, 0, time.UTC),
		},
		{
			name:    "Invalid format",
			input:   "not-a-date",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moment, err := gomoment.Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tt.input, err)
			}
			
			if !moment.Time().Equal(tt.wantTime) {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, moment.Time(), tt.wantTime)
			}
		})
	}
}

func TestTime(t *testing.T) {
	moment, _ := gomoment.NewMoment(testTime)
	
	if !moment.Time().Equal(testTime) {
		t.Errorf("Time() = %v, want %v", moment.Time(), testTime)
	}
}

// Benchmark tests to validate performance improvements
func BenchmarkFormat_Simple(b *testing.B) {
	moment, _ := gomoment.NewMoment(testTime)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		moment.Format("YYYY-MM-DD")
	}
}

func BenchmarkFormat_Complex(b *testing.B) {
	moment, _ := gomoment.NewMoment(testTime)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		moment.Format("YYYY-MM-DD HH:mm:ss Z")
	}
}

func BenchmarkNewMoment_CurrentTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoment.NewMoment()
	}
}

func BenchmarkNewMoment_ParseString(b *testing.B) {
	dateStr := "2023-05-15 14:30:45"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		gomoment.NewMoment(dateStr)
	}
}

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoment.Now()
	}
}

// Example tests for documentation
func ExampleNewMoment() {
	// Current time
	moment, _ := gomoment.NewMoment()
	formatted, _ := moment.Format("YYYY-MM-DD")
	_ = formatted // 2023-05-15 (or current date)
}

func ExampleNewMoment_withString() {
	moment, _ := gomoment.NewMoment("2023-05-15 14:30:45")
	formatted, _ := moment.Format("MM/DD/YYYY")
	_ = formatted // 05/15/2023
}

func ExampleMoment_Format() {
	moment, _ := gomoment.NewMoment("2023-05-15 14:30:45")
	
	iso, _ := moment.Format("YYYY-MM-DD")
	_ = iso // 2023-05-15
	
	us, _ := moment.Format("MM/DD/YYYY")
	_ = us // 05/15/2023
	
	time24, _ := moment.Format("HH:mm")
	_ = time24 // 14:30
}

func ExampleMoment_Must() {
	moment, _ := gomoment.NewMoment("2023-05-15")
	formatted := moment.Must("YYYY-MM-DD")
	_ = formatted // 2023-05-15
}