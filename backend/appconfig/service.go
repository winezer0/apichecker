package appconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"apichecker/backend/aiclient"
	"apichecker/backend/securestore"

	"gopkg.in/yaml.v3"
)

type PersistedConfig struct {
	KeysList    []KeyRow `yaml:"keys"`
	Prompt      string   `yaml:"prompt"`
	StartupPass string   `yaml:"startup_pass"`
}

// DefaultConfig 返回带默认值的配置对象。
func DefaultConfig() AppConfig {
	defaults := DefaultBaseURLs()
	return AppConfig{
		KeysList:  []KeyRow{},
		ChatState: ChatState{SelectedAPIIndex: -1},
		Defaults:  defaults,
		Prompt:    "你是一个智能助手，请用简洁明了的方式回答问题。",
	}
}

// DefaultBaseURLs 返回所有支持协议的默认 BaseURL。
func DefaultBaseURLs() map[string]string {
	return map[string]string{
		"openai":    aiclient.DefaultBaseURL("openai"),
		"anthropic": aiclient.DefaultBaseURL("anthropic"),
	}
}

// NormalizeConfig 对配置中的默认字段进行补全。
func NormalizeConfig(cfg AppConfig) AppConfig {
	defaults := DefaultBaseURLs()
	cfg.Defaults = defaults
	if cfg.ChatState.SelectedAPIIndex < -1 {
		cfg.ChatState.SelectedAPIIndex = -1
	}

	for index := range cfg.KeysList {
		cfg.KeysList[index].Protocol = strings.TrimSpace(cfg.KeysList[index].Protocol)
		cfg.KeysList[index].Alias = strings.TrimSpace(cfg.KeysList[index].Alias)
		cfg.KeysList[index].BaseURL = aiclient.NormalizeBaseURL(cfg.KeysList[index].Protocol, cfg.KeysList[index].BaseURL)
		cfg.KeysList[index].ModelCount = len(cfg.KeysList[index].Models)
	}

	return cfg
}

// Load 从 YAML 文件中加载配置，并使用启动密码或机器码解密。
func Load(path, machineID string, password string) (AppConfig, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}

	// 确定使用的解密密钥
	decryptionKey := machineID
	if password != "" {
		// 使用启动密码作为解密密钥
		decryptionKey = password
	}

	// 解密配置文件
	decryptedData, err := securestore.DecryptTextWithPassword(string(data), decryptionKey)
	if err != nil {
		return cfg, err
	}

	// 加载解密后的配置
	var persisted PersistedConfig
	if err = yaml.Unmarshal([]byte(decryptedData), &persisted); err != nil {
		return cfg, err
	}

	cfg.KeysList = persisted.KeysList
	cfg.Prompt = persisted.Prompt
	cfg.StartupPass = persisted.StartupPass

	return NormalizeConfig(cfg), nil
}

// Save 将配置写入 YAML 文件，并使用启动密码或机器码加密。
func Save(path, machineID string, cfg AppConfig) error {
	normalized := NormalizeConfig(cfg)
	persisted := PersistedConfig{
		KeysList:    normalized.KeysList,
		Prompt:      normalized.Prompt,
		StartupPass: normalized.StartupPass,
	}

	// 确定使用的加密密钥
	encryptionKey := machineID
	if normalized.StartupPass != "" {
		// 使用启动密码作为加密密钥
		encryptionKey = normalized.StartupPass
	}

	// 准备配置数据
	for i := range normalized.KeysList {
		normalized.KeysList[i].ModelCount = len(normalized.KeysList[i].Models)
	}
	persisted.KeysList = normalized.KeysList

	data, err := yaml.Marshal(&persisted)
	if err != nil {
		return err
	}

	// 加密整个配置文件
	encryptedData, err := securestore.EncryptTextWithPassword(string(data), encryptionKey)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(encryptedData), 0o600)
}

// ResolveConfigPath 根据应用名解析 YAML 配置文件路径。
func ResolveConfigPath(appName string) (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	configName := strings.TrimSpace(appName)
	if configName == "" {
		configName = strings.TrimSuffix(filepath.Base(executablePath), filepath.Ext(executablePath))
	}
	if filepath.Ext(configName) == "" {
		configName += ".yaml"
	}

	return filepath.Join(filepath.Dir(executablePath), configName), nil
}

// ResolveModelCachePath 根据应用名解析模型缓存文件路径。
func ResolveModelCachePath(appName string) (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	configName := strings.TrimSpace(appName)
	if configName == "" {
		configName = strings.TrimSuffix(filepath.Base(executablePath), filepath.Ext(executablePath))
	}

	return filepath.Join(filepath.Dir(executablePath), configName+".models.cache"), nil
}

// LoadModelCache 从 YAML 文件中加载模型缓存。
func LoadModelCache(path string) (ModelCache, error) {
	cfg := ModelCache{
		APIModels: make(map[string]map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// SaveModelCache 将模型缓存写入 YAML 文件。
func SaveModelCache(path string, cfg ModelCache) error {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
