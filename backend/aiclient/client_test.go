package aiclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIClient(t *testing.T) {
	// 模拟一个 OpenAI 服务
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": [{"id": "gpt-3.5-turbo"}]}`))
			return
		}
		if r.URL.Path == "/v1/chat/completions" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"choices": [{"message": {"content": "Test reply"}}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cfg := Config{
		Protocol: "openai",
		APIKey:   "test-key",
		BaseURL:  ts.URL + "/v1",
	}

	client := NewOpenAIClient(cfg)
	ctx := context.Background()

	// 测试验证
	valid, err := client.CheckValid(ctx)
	if err != nil || !valid {
		t.Errorf("Expected valid to be true, got err: %v", err)
	}

	// 测试获取模型
	models, err := client.GetModels(ctx)
	if err != nil || len(models) == 0 || models[0] != "gpt-3.5-turbo" {
		t.Errorf("Expected gpt-3.5-turbo, got %v (err: %v)", models, err)
	}

	// 测试对话
	reply, err := client.Chat(ctx, "gpt-3.5-turbo", "hello")
	if err != nil || reply != "Test reply" {
		t.Errorf("Expected 'Test reply', got '%s' (err: %v)", reply, err)
	}
}

func TestAnthropicClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/messages" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"content": [{"text": "Test anthropic reply"}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cfg := Config{
		Protocol: "anthropic",
		APIKey:   "test-key",
		BaseURL:  ts.URL + "/v1",
	}

	client := NewAnthropicClient(cfg)
	ctx := context.Background()

	// 测试验证
	valid, err := client.CheckValid(ctx)
	if err != nil || !valid {
		t.Errorf("Expected valid to be true, got err: %v", err)
	}

	// 测试对话
	reply, err := client.Chat(ctx, "claude-3", "hello")
	if err != nil || reply != "Test anthropic reply" {
		t.Errorf("Expected 'Test anthropic reply', got '%s' (err: %v)", reply, err)
	}
}
