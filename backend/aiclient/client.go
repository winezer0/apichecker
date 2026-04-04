package aiclient

import (
	"context"
	"errors"

	"apichecker/backend/machineid"
	"apichecker/backend/securestore"
)

// AIClient 接口定义了所有大模型客户端需要实现的方法
type AIClient interface {
	// CheckValid 验证 API Key 是否有效
	CheckValid(ctx context.Context) (bool, error)

	// GetModels 获取当前 API Key 支持的模型列表
	GetModels(ctx context.Context) ([]string, error)

	// Chat 发送一次聊天请求并获取回复
	Chat(ctx context.Context, model string, prompt string) (string, error)
}

// Config 定义了创建 AI 客户端的配置信息
type Config struct {
	Protocol string `json:"protocol"` // openai, anthropic, qwen
	APIKey   string `json:"apiKey"`
	BaseURL  string `json:"baseUrl"`
}

// getDecryptedAPIKey 如果 API Key 是加密的，则解密；否则返回原文
func getDecryptedAPIKey(apiKey string) (string, error) {
	// 检查是否是加密格式
	if len(apiKey) > 7 && apiKey[:7] == "aesgcm:" {
		machineID, err := machineid.GetMachineID()
		if err != nil {
			return "", err
		}
		return securestore.DecryptText(apiKey, machineID)
	}
	// 已经是明文，直接返回
	return apiKey, nil
}

// NewAIClient 根据配置工厂模式创建对应的 AIClient
func NewAIClient(cfg Config) (AIClient, error) {
	// 自动解密 API Key
	decryptedKey, err := getDecryptedAPIKey(cfg.APIKey)
	if err != nil {
		return nil, err
	}
	cfg.APIKey = decryptedKey

	switch cfg.Protocol {
	case "openai":
		return NewOpenAIClient(cfg), nil
	case "anthropic":
		return NewAnthropicClient(cfg), nil
	default:
		return nil, errors.New("unsupported protocol: " + cfg.Protocol)
	}
}
