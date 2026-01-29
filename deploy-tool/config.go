package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Build   BuildConfig          `json:"build"`
	Dev     map[string]DevConfig `json:"dev"`
	Stop    StopConfig           `json:"stop"`
	Timeout int                  `json:"timeout"` // 超时时间（秒）
}

type BuildConfig struct {
	Path []string `json:"path"`
}

type DevConfig struct {
	Command []string `json:"command"`
	Path    string   `json:"path"`
	Env     []string `json:"env"`
	Wait    []string `json:"wait"`
	Port    []int    `json:"port"` // 服务监听的端口列表
}

type StopConfig struct {
	Name []string `json:"name"`
}

var config Config

func LoadConfig() {
	if !FileExist("config.json") {
		panic("config not done")
	}

	configData, _ := os.ReadFile("config.json")
	json.Unmarshal(configData, &config)
}
