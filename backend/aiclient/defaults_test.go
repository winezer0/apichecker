package aiclient

import "testing"

func TestDefaultBaseURL(t *testing.T) {
	testCases := []struct {
		protocol string
		expected string
	}{
		{protocol: "openai", expected: "https://api.openai.com/v1"},
		{protocol: "anthropic", expected: "https://api.anthropic.com/v1"},
		{protocol: "unknown", expected: ""},
	}

	for _, testCase := range testCases {
		actual := DefaultBaseURL(testCase.protocol)
		if actual != testCase.expected {
			t.Fatalf("expected %s, got %s", testCase.expected, actual)
		}
	}
}

func TestNormalizeBaseURL(t *testing.T) {
	if actual := NormalizeBaseURL("openai", "  https://example.com/v1  "); actual != "https://example.com/v1" {
		t.Fatalf("expected custom base url, got %s", actual)
	}

	if actual := NormalizeBaseURL("anthropic", ""); actual != "https://api.anthropic.com/v1" {
		t.Fatalf("expected default anthropic url, got %s", actual)
	}
}
