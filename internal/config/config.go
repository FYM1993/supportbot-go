package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Redis    RedisConfig    `yaml:"redis"`
	DashScope DashScopeConfig `yaml:"dashscope"`
	Services ServicesConfig `yaml:"services"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// DashScopeConfig 通义千问配置
type DashScopeConfig struct {
	APIKey string `yaml:"apiKey"`
	Model  string `yaml:"model"`
}

// ServicesConfig 服务地址配置
type ServicesConfig struct {
	IMDemo              string `yaml:"imDemo"`
	QuestionClassifier  string `yaml:"questionClassifier"`
	Assistant           string `yaml:"assistant"`
	GeneralChat         string `yaml:"generalChat"`
	KnowledgeRAG        string `yaml:"knowledgeRag"`
	ProductService      string `yaml:"productService"`
	TradeService        string `yaml:"tradeService"`
	WorkOrderService    string `yaml:"workOrderService"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"` // debug, info, warn, error
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}

