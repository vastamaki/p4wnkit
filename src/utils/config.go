package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	Config     *AppConfig
	configPath string
	lock       sync.RWMutex
)

type AppConfig struct {
	ConfigDir      string
	HtbAccessToken string `json:"htb_access_token"`
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving home directory: %v\n", err)
		os.Exit(1)
	}

	configPath = filepath.Join(homeDir, ".p4wnkit", "config.json")
}

func InitializeConfig() (config *AppConfig) {
	dirPath := filepath.Dir(configPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Join(dirPath, "openvpn"), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Default config file created:", configPath)
	}

	LoadConfig()

	return Config
}

func createDefaultConfig() error {
	defaultConfig := AppConfig{
		HtbAccessToken: "",
	}

	return saveConfig(configPath, defaultConfig)
}

func LoadConfig() {
	lock.Lock()
	defer lock.Unlock()

	dirPath := filepath.Dir(configPath)

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(configFile, &Config); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling config file: %v\n", err)
		os.Exit(1)
	}

	Config.ConfigDir = dirPath
}

func saveConfig(filePath string, cfg AppConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	return encoder.Encode(cfg)
}

func SaveConfig(cfg *AppConfig) error {
	lock.Lock()
	defer lock.Unlock()
	return saveConfig(configPath, *cfg)
}

func GetConfig() AppConfig {
	lock.RLock()
	defer lock.RUnlock()

	return *Config
}
