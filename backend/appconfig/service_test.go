package appconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"apichecker/backend/securestore"
)

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "apichecker.yaml")
	machineID := "test-machine-id"

	cfg := AppConfig{
		KeysList: []KeyRow{
			{
				Protocol: "openai",
				Alias:    "default-openai",
				APIKey:   "sk-openai-secret",
				BaseURL:  "",
				Success:  true,
				Models:   []string{"gpt-4o-mini", "gpt-4.1"},
			},
			{
				Protocol: "anthropic",
				Alias:    "default-claude",
				APIKey:   "sk-ant-secret",
				BaseURL:  "",
				Success:  true,
				Models:   []string{"claude-3-5-sonnet-20241022"},
			},
		},
		ChatState: ChatState{
			SelectedAPIIndex: 1,
			SelectedModel:    "claude-3-5-sonnet-20241022",
		},
	}

	if err := Save(configPath, machineID, cfg); err != nil {
		t.Fatalf("save config failed: %v", err)
	}

	rawContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config file failed: %v", err)
	}
	if strings.Contains(string(rawContent), "sk-openai-secret") || strings.Contains(string(rawContent), "sk-ant-secret") {
		t.Fatal("expected api key to be encrypted in yaml file")
	}

	loaded, err := Load(configPath, machineID)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if !strings.HasPrefix(loaded.KeysList[0].APIKey, "aesgcm:") {
		t.Fatalf("expected openai key to be encrypted, got %s", loaded.KeysList[0].APIKey)
	}

	decryptedOpenAI, err := securestore.DecryptText(loaded.KeysList[0].APIKey, machineID)
	if err != nil {
		t.Fatalf("decrypt openai key failed: %v", err)
	}
	if decryptedOpenAI != "sk-openai-secret" {
		t.Fatalf("expected decrypted openai key to be sk-openai-secret, got %s", decryptedOpenAI)
	}

	if loaded.KeysList[0].BaseURL != "https://api.openai.com/v1" {
		t.Fatalf("expected default openai base url, got %s", loaded.KeysList[0].BaseURL)
	}
	if loaded.KeysList[1].BaseURL != "https://api.anthropic.com/v1" {
		t.Fatalf("expected default anthropic base url, got %s", loaded.KeysList[1].BaseURL)
	}
	if loaded.KeysList[0].ModelCount != 2 {
		t.Fatalf("expected model count 2, got %d", loaded.KeysList[0].ModelCount)
	}
}

func TestLoadMissingFileReturnsDefaultConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "missing.yaml")

	_, err := Load(configPath, "test-machine-id")
	if err != nil {
		t.Fatalf("expected missing file to return default config, got error: %v", err)
	}
}
