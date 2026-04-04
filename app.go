package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"apichecker/backend/aiclient"
	"apichecker/backend/appconfig"
	"apichecker/backend/cache"
	"apichecker/backend/checker"
	"apichecker/backend/machineid"
	"apichecker/backend/securestore"
)

const appConfigFileName = "apichecker.yaml"

// App struct
type App struct {
	ctx           context.Context
	configPath    string
	machineID     string
	config        appconfig.AppConfig
	configLoadErr error
	mu            sync.RWMutex
}

// NewApp 创建新的应用实例。
func NewApp() *App {
	return &App{
		config: appconfig.DefaultConfig(),
	}
}

// startup 在应用启动时初始化上下文和本地配置。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	configPath, err := appconfig.ResolveConfigPath(appConfigFileName)
	if err != nil {
		a.configLoadErr = err
		return
	}

	machineCode, err := machineid.GetMachineID()
	if err != nil {
		a.configPath = configPath
		a.configLoadErr = err
		return
	}

	loadedConfig, err := appconfig.Load(configPath, machineCode, "")
	a.configPath = configPath
	a.machineID = machineCode
	a.config = loadedConfig
	a.configLoadErr = err
}

// GetStartupPassStatus 获取启动密码状态
func (a *App) GetStartupPassStatus() (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.configLoadErr != nil {
		return false, a.configLoadErr
	}

	return a.config.StartupPass != "", nil
}

// shutdown 在应用关闭时执行最后一次配置落盘。
func (a *App) shutdown(ctx context.Context) {
	a.mu.RLock()
	cfg := a.config
	configPath := a.configPath
	machineCode := a.machineID
	a.mu.RUnlock()

	if configPath == "" || machineCode == "" {
		return
	}

	_ = appconfig.Save(configPath, machineCode, cfg)
}

// GetProtocols 获取支持的协议列表
func (a *App) GetProtocols() []string {
	return []string{"openai", "anthropic"}
}

// GetDefaultBaseURLs 获取各协议对应的默认 BaseURL。
func (a *App) GetDefaultBaseURLs() map[string]string {
	return appconfig.DefaultBaseURLs()
}

// GetAppConfig 获取当前已加载的应用配置（加密状态返回给前端）。
func (a *App) GetAppConfig() (appconfig.AppConfig, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.configLoadErr != nil {
		return appconfig.DefaultConfig(), a.configLoadErr
	}

	return a.config, nil
}

// SaveAppConfig 保存当前应用配置到 YAML 文件并返回更新后的配置。
func (a *App) SaveAppConfig(cfg appconfig.AppConfig) (appconfig.AppConfig, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return appconfig.DefaultConfig(), errors.New("config path is not initialized")
	}
	if a.machineID == "" {
		if a.configLoadErr != nil {
			return appconfig.DefaultConfig(), a.configLoadErr
		}
		return appconfig.DefaultConfig(), errors.New("machine id is not initialized")
	}

	// 保留现有的启动密码
	cfg.StartupPass = a.config.StartupPass

	normalized := appconfig.NormalizeConfig(cfg)
	if err := appconfig.Save(a.configPath, a.machineID, normalized); err != nil {
		return appconfig.DefaultConfig(), err
	}

	loadedConfig, err := appconfig.Load(a.configPath, a.machineID, a.config.StartupPass)
	if err != nil {
		return appconfig.DefaultConfig(), err
	}

	a.config = loadedConfig
	a.configLoadErr = nil
	return loadedConfig, nil
}

// CheckKey 验证单个API Key，获取模型
func (a *App) CheckKey(protocol, alias, apiKey, baseURL string) checker.CheckResult {
	return checker.CheckKey(protocol, alias, apiKey, baseURL)
}

// BatchCheckKeys 批量验证API Key
func (a *App) BatchCheckKeys(protocol string, aliases []string, keys []string, baseURL string) checker.BatchCheckResult {
	return checker.BatchCheckKeys(protocol, aliases, keys, baseURL)
}

// Chat 聊天功能，用于调试窗口
func (a *App) Chat(protocol, apiKey, baseURL, model, prompt string) (string, error) {
	cfg := aiclient.Config{
		Protocol: protocol,
		APIKey:   apiKey,
		BaseURL:  aiclient.NormalizeBaseURL(protocol, baseURL),
	}
	client, err := aiclient.NewAIClient(cfg)
	if err != nil {
		return "", err
	}
	return client.Chat(context.Background(), model, prompt)
}

// TestModel 测试单个模型的智能度并缓存结果
func (a *App) TestModel(protocol, apiKey, baseURL, model, prompt string) (string, error) {
	reply, err := a.Chat(protocol, apiKey, baseURL, model, prompt)
	if err == nil {
		cacheKey := cache.GenerateCacheKey(protocol, apiKey, model)
		cache.GlobalCache.Set(cacheKey, reply)
	}
	return reply, err
}

// GetCachedModelResult 获取模型智能度测试的缓存结果
func (a *App) GetCachedModelResult(protocol, apiKey, model string) string {
	cacheKey := cache.GenerateCacheKey(protocol, apiKey, model)
	if val, ok := cache.GlobalCache.Get(cacheKey); ok {
		return val
	}
	return ""
}

// ExportEncryptedList 导出加密的列表数据
func (a *App) ExportEncryptedList() ([]map[string]string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.configLoadErr != nil {
		return nil, a.configLoadErr
	}

	result := make([]map[string]string, len(a.config.KeysList))
	for i, row := range a.config.KeysList {
		result[i] = map[string]string{
			"protocol":   row.Protocol,
			"alias":      row.Alias,
			"apiKey":     row.APIKey,
			"baseUrl":    row.BaseURL,
			"success":    fmt.Sprintf("%t", row.Success),
			"modelCount": fmt.Sprintf("%d", row.ModelCount),
			"errorMsg":   row.ErrorMsg,
		}
	}
	return result, nil
}

// DecryptAPIKey 解密单个 API Key，仅在前端需要时调用
func (a *App) DecryptAPIKey(encryptedKey string) (string, error) {
	a.mu.RLock()
	machineCode := a.machineID
	a.mu.RUnlock()

	if machineCode == "" {
		return "", errors.New("machine id is not initialized")
	}

	return securestore.DecryptText(encryptedKey, machineCode)
}

// VerifyStartupPass 验证启动密码
func (a *App) VerifyStartupPass(password string) (bool, error) {
	a.mu.RLock()
	machineCode := a.machineID
	storedPass := a.config.StartupPass
	a.mu.RUnlock()

	if machineCode == "" {
		return false, errors.New("machine id is not initialized")
	}

	if storedPass == "" {
		return true, nil
	}

	decrypted, err := securestore.DecryptText(storedPass, machineCode)
	if err != nil {
		return false, err
	}

	return decrypted == password, nil
}

// GetPrompt 获取智能提示词
func (a *App) GetPrompt() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.config.Prompt
}

// SetPrompt 设置智能提示词
func (a *App) SetPrompt(prompt string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return errors.New("config path is not initialized")
	}
	if a.machineID == "" {
		if a.configLoadErr != nil {
			return a.configLoadErr
		}
		return errors.New("machine id is not initialized")
	}

	a.config.Prompt = prompt
	normalized := appconfig.NormalizeConfig(a.config)
	if err := appconfig.Save(a.configPath, a.machineID, normalized); err != nil {
		return err
	}

	loadedConfig, err := appconfig.Load(a.configPath, a.machineID, a.config.StartupPass)
	if err != nil {
		return err
	}

	a.config = loadedConfig
	return nil
}

// SetStartupPass 设置启动密码
func (a *App) SetStartupPass(password string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return errors.New("config path is not initialized")
	}
	if a.machineID == "" {
		if a.configLoadErr != nil {
			return a.configLoadErr
		}
		return errors.New("machine id is not initialized")
	}

	a.config.StartupPass = password
	normalized := appconfig.NormalizeConfig(a.config)
	if err := appconfig.Save(a.configPath, a.machineID, normalized); err != nil {
		return err
	}

	loadedConfig, err := appconfig.Load(a.configPath, a.machineID, a.config.StartupPass)
	if err != nil {
		return err
	}

	a.config = loadedConfig
	return nil
}

// GetModelCache 获取模型缓存
func (a *App) GetModelCache() (appconfig.ModelCache, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return appconfig.ModelCache{}, errors.New("config path is not initialized")
	}

	modelCachePath, err := appconfig.ResolveModelCachePath("apichecker")
	if err != nil {
		return appconfig.ModelCache{}, err
	}

	return appconfig.LoadModelCache(modelCachePath)
}

// SaveModelCache 保存模型缓存
func (a *App) SaveModelCache(cache appconfig.ModelCache) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return errors.New("config path is not initialized")
	}

	modelCachePath, err := appconfig.ResolveModelCachePath("apichecker")
	if err != nil {
		return err
	}

	return appconfig.SaveModelCache(modelCachePath, cache)
}

// LoadConfigWithPassword 使用启动密码加载配置文件
func (a *App) LoadConfigWithPassword(password string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.machineID == "" {
		return errors.New("machine id is not initialized")
	}
	if a.configPath == "" {
		return errors.New("config path is not initialized")
	}

	// 使用密码或机器码加载配置文件
	loadedConfig, err := appconfig.Load(a.configPath, a.machineID, password)
	if err != nil {
		return err
	}

	// 更新配置
	a.config = loadedConfig
	a.configLoadErr = nil

	return nil
}

// ClearStartupPass 清除启动密码
func (a *App) ClearStartupPass() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.configPath == "" {
		return errors.New("config path is not initialized")
	}
	if a.machineID == "" {
		if a.configLoadErr != nil {
			return a.configLoadErr
		}
		return errors.New("machine id is not initialized")
	}

	// 直接清除启动密码
	a.config.StartupPass = ""

	normalized := appconfig.NormalizeConfig(a.config)
	if err := appconfig.Save(a.configPath, a.machineID, normalized); err != nil {
		return err
	}

	// 直接使用当前的配置对象
	a.config = normalized
	return nil
}
