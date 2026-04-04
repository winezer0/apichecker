package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// AnthropicClient Anthropic协议兼容客户端
type AnthropicClient struct {
	Config Config
	client *http.Client
}

// NewAnthropicClient 创建一个新的Anthropic兼容客户端
func NewAnthropicClient(cfg Config) *AnthropicClient {
	cfg.BaseURL = NormalizeBaseURL(cfg.Protocol, cfg.BaseURL)
	return &AnthropicClient{
		Config: cfg,
		client: &http.Client{},
	}
}

func (c *AnthropicClient) doRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.Config.BaseURL, path)
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", c.Config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

// CheckValid 验证API Key是否有效
func (c *AnthropicClient) CheckValid(ctx context.Context) (bool, error) {
	// Anthropic没有通用的/models接口用于列出模型，通常用无效的请求来验证身份
	// 我们构造一个简单的messages请求，只看是否返回401来判断
	reqBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307",
		"max_tokens": 1,
		"messages": []map[string]string{
			{"role": "user", "content": "ping"},
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}
	resp, err := c.doRequest(ctx, "POST", "/messages", b)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	// 如果是401/403，则是无效的key，200则是有效
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return false, fmt.Errorf("invalid api key: %d", resp.StatusCode)
	}
	return true, nil
}

// GetModels 获取当前API Key支持的模型列表
func (c *AnthropicClient) GetModels(ctx context.Context) ([]string, error) {
	// Anthropic官方目前不提供List Models API，所以我们可以返回固定的支持列表或通过尝试调用来判断
	// 为了简单起见，返回Claude3系列模型
	return []string{
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}, nil
}

// Chat 发送一次聊天请求并获取回复
func (c *AnthropicClient) Chat(ctx context.Context, model string, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":      model,
		"max_tokens": 1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := c.doRequest(ctx, "POST", "/messages", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("invalid status code: %d, response: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", errors.New("no choice returned")
}
