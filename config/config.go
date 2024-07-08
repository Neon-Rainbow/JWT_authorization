package config

import (
	"encoding/json"
	"log"
	"os"
)

var appConfig *Config

type MySQL struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type JWT struct {
	Secret string `json:"secret"`
}

type Config struct {
	Address        string `json:"address"`
	Port           int    `json:"port"`
	MySQL          MySQL  `json:"MySQL"`
	JWT            JWT    `json:"JWT"`
	PasswordSecret string `json:"passwordSecret"`
}

// LoadConfig is a function to load config from file
func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}(file)

	appConfig = &Config{}
	if err := json.NewDecoder(file).Decode(appConfig); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
		return err
	}
	return nil
}

// GetConfig is a function to get config
func GetConfig() *Config {
	return appConfig
}
