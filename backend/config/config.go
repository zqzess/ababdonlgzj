package config

import (
	"abandonlgzj/logger"
	"abandonlgzj/types"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	jsonFilePath *string
	Config       types.Conf
	isFileConfig bool // 是否成功加载到外部配置文件
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (s *Config) Start(path *string) {
	s.jsonFilePath = path
	_, err := os.Stat(*s.jsonFilePath)
	if err != nil {
		logger.Console.Error("%s  no such file or directory", *s.jsonFilePath)
		s.isFileConfig = false
		return
	}
	s.GetConfig()
}

// GetConfig 获取配置文件
func (s *Config) GetConfig() {
	byteValue, err := os.ReadFile(*s.jsonFilePath)
	if err != nil {
		logger.Console.Error("Error reading JSON file:", err)
		//os.Exit(1)
		return
	}

	err2 := json.Unmarshal(byteValue, &s.Config)
	if err != nil {
		logger.Console.Error("Error parsing JSON data:", err2)
		//os.Exit(1)
		return
	}
	s.isFileConfig = true
	logger.Console.Info("config load success!")
	logger.Console.Debug(s.Config.Welcome)
}

// GetBackPort 获取端口
func (s *Config) GetBackPort() int {
	if s.isFileConfig {
		return s.Config.BackPort
	} else {
		return 9091
	}
}

// GetFrontPort 获取端口
func (s *Config) GetFrontPort() int {
	if s.isFileConfig {
		return s.Config.FrontPort
	} else {
		return 9091
	}
}

func (s *Config) GetUploadSize() int64 {
	if s.isFileConfig {
		return s.Config.UploadSize
	} else {
		return 100
	}
}

func (s *Config) GetTmpPath() string {
	if s.isFileConfig {
		if s.Config.Tmp != "null" && s.Config.Tmp != "" {
			return s.Config.Tmp
		}
	}
	exePath, err := os.Executable()
	if err != nil {
		logger.Console.Error(err.Error())
	}
	dirPath := filepath.Dir(exePath)
	return dirPath + "/tmpFile"
}

func (s *Config) GetRemoveTime() int {
	if s.isFileConfig {
		return s.Config.AutoRemoveTime
	} else {
		return 30
	}
}

func (s *Config) ReloadConfig() {
	s.GetConfig()
}
