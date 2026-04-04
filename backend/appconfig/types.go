package appconfig

// FormState 表示顶部新增表单的当前状态。
type FormState struct {
	Protocol string `json:"protocol" yaml:"protocol"`
	Alias    string `json:"alias" yaml:"alias"`
	APIKey   string `json:"apiKey" yaml:"apiKey"`
	BaseURL  string `json:"baseUrl" yaml:"base_url"`
}

// KeyRow 表示单个 API 配置项。
type KeyRow struct {
	Protocol   string   `json:"protocol" yaml:"protocol"`
	Alias      string   `json:"alias" yaml:"alias"`
	APIKey     string   `json:"apiKey" yaml:"apiKey"`
	BaseURL    string   `json:"baseUrl" yaml:"base_url"`
	Success    bool     `json:"success" yaml:"success"`
	Models     []string `json:"models" yaml:"models"`
	ModelCount int      `json:"modelCount" yaml:"model_count"`
	ErrorMsg   string   `json:"errorMsg" yaml:"error_msg"`
}

// ModelCache 表示模型测试缓存。
type ModelCache struct {
	APIModels map[string]map[string]string `json:"apiModels" yaml:"api_models"` // api类型@api别名 -> 模型名称 -> 模型响应
}

// ChatState 表示底部聊天区域的选择状态。
type ChatState struct {
	SelectedAPIIndex int    `json:"selectedApiIndex" yaml:"selected_api_index"`
	SelectedModel    string `json:"selectedModel" yaml:"selected_model"`
}

// AppConfig 表示应用完整配置。
type AppConfig struct {
	KeysList    []KeyRow          `json:"keysList" yaml:"keys"`
	ChatState   ChatState         `json:"chatState" yaml:"chat_state"`
	Defaults    map[string]string `json:"defaults" yaml:"defaults"`
	StartupPass string            `json:"startupPass" yaml:"startup_pass"`
	Prompt      string            `json:"prompt" yaml:"prompt"`
}
