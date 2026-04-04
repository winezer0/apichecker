package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": [{"id": "gpt-3.5-turbo"}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	result := CheckKey("openai", "test-alias", "test-key", ts.URL+"/v1")

	if !result.IsValid {
		t.Errorf("Expected valid, got invalid: %s", result.ErrorMsg)
	}

	if result.Alias != "test-alias" {
		t.Errorf("Expected alias test-alias, got %s", result.Alias)
	}

	if len(result.Models) == 0 || result.Models[0] != "gpt-3.5-turbo" {
		t.Errorf("Expected model gpt-3.5-turbo, got %v", result.Models)
	}
}

func TestBatchCheckKeys(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.URL.Path == "/v1/models" {
			w.Write([]byte(`{"data": [{"id": "gpt-3.5-turbo"}]}`))
		} else {
			w.Write([]byte(`{"choices": [{"message": {"content": "ok"}}]}`))
		}
	}))
	defer ts.Close()

	keys := []string{"key1", "key2", "key3"}
	aliases := []string{"alias1", "alias2", "alias3"}
	res := BatchCheckKeys("openai", aliases, keys, ts.URL+"/v1")

	if len(res.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(res.Results))
	}
}
