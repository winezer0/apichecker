package aiclient

import "strings"

// DefaultBaseURL 返回协议对应的默认 BaseURL。
func DefaultBaseURL(protocol string) string {
	switch strings.ToLower(strings.TrimSpace(protocol)) {
	case "anthropic":
		return "https://api.anthropic.com/v1"
	case "openai":
		return "https://api.openai.com/v1"
	default:
		return ""
	}
}

// NormalizeBaseURL 归一化 BaseURL，优先使用用户输入，否则回退到协议默认值。
func NormalizeBaseURL(protocol, baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed != "" {
		return trimmed
	}
	return DefaultBaseURL(protocol)
}
