package render

import (
	"testing"
)

func TestConvertFormat(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"dd-MM-yyyy", "02-01-2006"},
		{"yyyy/MM/dd", "2006/01/02"},
		{"HH:mm:ss", "15:04:05"},
		{"dd MMMM yyyy HH:mm", "02 MMMM 2006 15:04"},
		{"dd-MM-yy", "02-01-yy"},
		{"\\d{4}-\\d{2}", "\\d{4}-\\d{2}"},
		{"", ""},
		{"HH\\:mm", "15:04"},
	}
	for _, tt := range tests {
		got := convertFormat(tt.input)
		if got != tt.want {
			t.Errorf("convertFormat(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestApplyFormatCustomSpecifiers(t *testing.T) {
	tests := []struct {
		val    string
		fmtStr string
		want   string
	}{
		{"2026-06-24", "dd-MM-yyyy", "24-06-2026"},
		{"2026-06-24", "yyyy/MM/dd", "2026/06/24"},
		{"2026-06-24T15:04:05Z", "HH:mm:ss", "15:04:05"},
	}
	for _, tt := range tests {
		got := applyFormat(tt.val, tt.fmtStr)
		if got != tt.want {
			t.Errorf("applyFormat(%q, %q) = %q, want %q", tt.val, tt.fmtStr, got, tt.want)
		}
	}
}

func TestApplyFormatGoLayout(t *testing.T) {
	tests := []struct {
		val    string
		fmtStr string
		want   string
	}{
		{"2026-06-24", "02-01-2006", "24-06-2026"},
	}
	for _, tt := range tests {
		got := applyFormat(tt.val, tt.fmtStr)
		if got != tt.want {
			t.Errorf("applyFormat(%q, %q) = %q, want %q", tt.val, tt.fmtStr, got, tt.want)
		}
	}
}

func TestApplyFormatUnmatched(t *testing.T) {
	got := applyFormat("not-a-date", "dd-MM-yyyy")
	if got != "not-a-date" {
		t.Errorf("applyFormat should return original value when parsing fails, got %q", got)
	}
}
