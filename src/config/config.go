package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 应用配置结构
type Config struct {
	Persons    []string `json:"persons"`     // 站会主持人列表
	FirstData  int64    `json:"first_data"`  // 第一个站会的日期，格式为yyyyMMdd文本转数字
	FirstIndex int      `json:"first_index"` // 第一个站会的主持人在persons列表中的索引，从0开始
	SlackUrl   string   `json:"slack_url"`   // Slack 通知 URL
	DataFile   string   `json:"data_file"`   // 通知数据文件名
}

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	// 确保文件路径是绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
