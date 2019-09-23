package base62

import "testing"

func TestEncode(t *testing.T) {
	testData := []struct {
		key      string
		expected uint64
	}{
		{"0", 0},
		{"a", 10},
		{"aa", 630},
		{"abc123EFG", 2222821365901088},
		{"hjNv8tS3K", 3781504209452600},
	}
	for _, testCase := range testData {
		n, err := Decode(testCase.key)
		if n != testCase.expected {
			t.Fatalf("Expected %v, but got %v", testCase.expected, n)
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestDecode(t *testing.T) {
	testData := []struct {
		key      string
		expected uint64
	}{
		{"0", 0},
		{"a", 10},
		{"aa", 630},
		{"abc123EFG", 2222821365901088},
		{"hjNv8tS3K", 3781504209452600},
	}
	for _, testCase := range testData {
		n, err := Decode(testCase.key)
		if n != testCase.expected {
			t.Fatalf("expected %v, but got %v", testCase.expected, n)
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
