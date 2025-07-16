package unpackage

import (
	"testing"
)

func TestUnpackage(t *testing.T) {
	tests := []struct {
		input    string
		want     string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
	}

	for _, tt := range tests {
		got, err := Unpackage(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Unpackage(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("Unpackage(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}