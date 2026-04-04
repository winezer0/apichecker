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

// OpenAIClient OpenAI协议兼容的客户端实现
type OpenAIClient struct {
	Config Config
	client *http.Client
}

// NewOpenAIClient 创建一个新的OpenAI兼容客户端
func NewOpenAIClient(cfg Config) *OpenAIClient {
	cfg.BaseURL = NormalizeBaseURL(cfg.Protocol, cfg.BaseURL)
	return &OpenAIClient{
		Config: cfg,
		client: &http.Client{},
	}
}

func (c *OpenAIClient) doRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
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
	req.Header.Set("Authorization", "Bearer "+c.Config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

// CheckValid 验证API Key是否有效
func (c *OpenAIClient) CheckValid(ctx context.Context) (bool, error) {
	resp, err := c.doRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("invalid status code: %d", resp.StatusCode)
}

// GetModels 获取当前API Key支持的模型列表
func (c *OpenAIClient) GetModels(ctx context.Context) ([]string, error) {
	resp, err := c.doRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}
	return models, nil
}

// Chat 发送一次聊天请求并获取回复
func (c *OpenAIClient) Chat(ctx context.Context, model string, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := c.doRequest(ctx, "POST", "/chat/completions", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("invalid status code: %d, response: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", errors.New("no choice returned")
}
