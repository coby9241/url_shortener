package utils

import (
	"testing"
)

func TestEncodeToBase62(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "0"},
		{"\x01", "1"},
		{"\x3d", "Z"},
		{"\x3e", "10"},
		{"\x7b", "1Z"},
	}

	for _, test := range tests {
		result := EncodeToBase62(test.input)
		if result != test.expected {
			t.Errorf("EncodeToBase62(%q) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestGenerateURLHash(t *testing.T) {
	url1 := "https://example.com"
	url2 := "https://example.com/some/path"

	hash1 := GenerateURLHash(url1, 7)
	hash2 := GenerateURLHash(url2, 8)

	if len(hash1) != 7 {
		t.Errorf("expected hash length to be 7, got %d", len(hash1))
	}

	if len(hash2) != 8 {
		t.Errorf("expected hash length to be 8, got %d", len(hash2))
	}

	if hash1 == hash2 {
		t.Errorf("expected hashes to be different for different URLs, got identical: %s", hash1)
	}

	// Test determinism
	hash1Verify := GenerateURLHash(url1, 7)
	if hash1 != hash1Verify {
		t.Errorf("expected hash to be deterministic, got %s then %s", hash1, hash1Verify)
	}
}
