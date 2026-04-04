package checker

import (
	"context"
	"strings"
	"sync"
	"time"

	"apichecker/backend/aiclient"
)

// CheckResult 包含单次 API Key 检测的结果
type CheckResult struct {
	Alias    string   `json:"alias"`
	Key      string   `json:"key"`
	Protocol string   `json:"protocol"`
	IsValid  bool     `json:"isValid"`
	Models   []string `json:"models"`
	ErrorMsg string   `json:"errorMsg"`
}

// CheckKey 验证单个API Key，获取模型
func CheckKey(protocol, alias, apiKey, baseURL string) CheckResult {
	result := CheckResult{
		Alias:    alias,
		Key:      apiKey,
		Protocol: protocol,
	}

	cfg := aiclient.Config{
		Protocol: protocol,
		APIKey:   apiKey,
		BaseURL:  baseURL,
	}

	client, err := aiclient.NewAIClient(cfg)
	if err != nil {
		result.ErrorMsg = "Client initialization failed: " + err.Error()
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. 验证有效性
	isValid, err := client.CheckValid(ctx)
	if err != nil || !isValid {
		result.IsValid = false
		if err != nil {
			result.ErrorMsg = "Validation failed: " + err.Error()
		} else {
			result.ErrorMsg = "Invalid API Key"
		}
		return result
	}
	result.IsValid = true

	// 2. 获取模型
	models, err := client.GetModels(ctx)
	if err != nil {
		// 某些接口可能不支持获取模型，只记录错误但不算无效
		result.ErrorMsg = "Failed to get models: " + err.Error()
	} else {
		result.Models = models
	}

	return result
}

// BatchCheckResult 批量验证结果
type BatchCheckResult struct {
	Results []CheckResult `json:"results"`
}

// BatchCheckKeys 并发验证多个API Key
func BatchCheckKeys(protocol string, aliases []string, keys []string, baseURL string) BatchCheckResult {
	var wg sync.WaitGroup
	resultChan := make(chan CheckResult, len(keys))

	for i, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}

		alias := ""
		if i < len(aliases) {
			alias = aliases[i]
		}

		wg.Add(1)
		go func(key, aliasName string) {
			defer wg.Done()
			res := CheckKey(protocol, aliasName, key, baseURL)
			resultChan <- res
		}(k, alias)
	}

	wg.Wait()
	close(resultChan)

	var finalResult BatchCheckResult
	for r := range resultChan {
		finalResult.Results = append(finalResult.Results, r)
	}

	return finalResult
}
