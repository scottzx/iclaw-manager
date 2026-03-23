package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"iclaw-admin-api/internal/model"
)

// ConfigPaths 配置路径
var ConfigPaths = struct {
	OpenClawJSON string
	EnvFile     string
}{
	OpenClawJSON: GetConfigDir() + "/openclaw.json",
	EnvFile:     GetConfigDir() + "/env",
}

// GetConfig 获取完整配置
func GetConfig() (map[string]interface{}, error) {
	data, err := os.ReadFile(ConfigPaths.OpenClawJSON)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SaveConfig 保存配置
func SaveConfig(config map[string]interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(ConfigPaths.OpenClawJSON)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(ConfigPaths.OpenClawJSON, data, 0600)
}

// GetEnvValue 获取环境变量
func GetEnvValue(key string) (string, error) {
	data, err := os.ReadFile(ConfigPaths.EnvFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export ") {
			line = line[7:]
		}
		if strings.HasPrefix(line, key+"=") {
			value := strings.TrimPrefix(line, key+"=")
			value = strings.Trim(strings.Trim(value, "\"'"), "\"")
			return value, nil
		}
	}

	return "", nil
}

// SaveEnvValue 保存环境变量
func SaveEnvValue(key, value string) error {
	data, err := os.ReadFile(ConfigPaths.EnvFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := []string{}
	if len(data) > 0 {
		lines = strings.Split(string(data), "\n")
	}

	// 更新或添加
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "export ") {
			trimmed = trimmed[7:]
		}
		if strings.HasPrefix(trimmed, key+"=") {
			lines[i] = fmt.Sprintf("export %s=\"%s\"", key, value)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("export %s=\"%s\"", key, value))
	}

	// 确保目录存在
	dir := filepath.Dir(ConfigPaths.EnvFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(ConfigPaths.EnvFile, []byte(strings.Join(lines, "\n")+"\n"), 0600)
}

// GetAIConfig 获取 AI 配置概览
func GetAIConfig() (*model.AIConfigOverview, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	overview := &model.AIConfigOverview{
		ConfiguredProviders: []model.ConfiguredProvider{},
		AvailableModels:    []string{},
	}

	// 获取主模型
	if agents, ok := config["agents"].(map[string]interface{}); ok {
		if defaults, ok := agents["defaults"].(map[string]interface{}); ok {
			if model, ok := defaults["model"].(map[string]interface{}); ok {
				if primary, ok := model["primary"].(string); ok {
					overview.PrimaryModel = &primary
				}
			}
		}
	}

	// 获取配置的 providers
	if models, ok := config["models"].(map[string]interface{}); ok {
		if providers, ok := models["providers"].(map[string]interface{}); ok {
			for name, pData := range providers {
				pMap, ok := pData.(map[string]interface{})
				if !ok {
					continue
				}

				cp := model.ConfiguredProvider{
					Name:   name,
					Models: []model.ConfiguredModel{},
				}

				if baseURL, ok := pMap["baseUrl"].(string); ok {
					cp.BaseURL = baseURL
				}

				if apiKey, ok := pMap["apiKey"].(string); ok && apiKey != "" {
					cp.HasAPIKey = true
					masked := apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
					cp.APIKeyMasked = &masked
				}

				if modelsList, ok := pMap["models"].([]interface{}); ok {
					for _, mData := range modelsList {
						if mMap, ok := mData.(map[string]interface{}); ok {
							cm := model.ConfiguredModel{}
							if id, ok := mMap["id"].(string); ok {
								cm.ID = id
								cm.FullID = name + "/" + id
							}
							if name, ok := mMap["name"].(string); ok {
								cm.Name = name
							}
							if api, ok := mMap["api"].(string); ok {
								cm.APIType = &api
							}
							if cw, ok := mMap["contextWindow"].(float64); ok {
								cwUint := uint32(cw)
								cm.ContextWindow = &cwUint
							}
							if mt, ok := mMap["maxTokens"].(float64); ok {
								mtUint := uint32(mt)
								cm.MaxTokens = &mtUint
							}
							if overview.PrimaryModel != nil && cm.FullID == *overview.PrimaryModel {
								cm.IsPrimary = true
							}
							cp.Models = append(cp.Models, cm)
						}
					}
				}

				overview.ConfiguredProviders = append(overview.ConfiguredProviders, cp)
			}
		}
	}

	// 获取可用模型列表
	if agents, ok := config["agents"].(map[string]interface{}); ok {
		if defaults, ok := agents["defaults"].(map[string]interface{}); ok {
			if modelsMap, ok := defaults["models"].(map[string]interface{}); ok {
				for fullID := range modelsMap {
					overview.AvailableModels = append(overview.AvailableModels, fullID)
				}
			}
		}
	}

	return overview, nil
}

// ValidateConfig 验证配置
func ValidateConfig(config map[string]interface{}) (bool, []string) {
	var errors []string

	// 基本 JSON 验证
	data, err := json.Marshal(config)
	if err != nil {
		errors = append(errors, fmt.Sprintf("无效的 JSON: %v", err))
		return false, errors
	}

	// 重新解析验证
	var test map[string]interface{}
	if err := json.Unmarshal(data, &test); err != nil {
		errors = append(errors, fmt.Sprintf("配置格式错误: %v", err))
		return false, errors
	}

	return len(errors) == 0, errors
}

// GetGatewayToken 获取或创建 gateway token
func GetGatewayToken() (string, error) {
	config, err := GetConfig()
	if err != nil {
		return "", err
	}

	// 尝试从配置中获取
	if gateway, ok := config["gateway"].(map[string]interface{}); ok {
		if auth, ok := gateway["auth"].(map[string]interface{}); ok {
			if token, ok := auth["token"].(string); ok && token != "" {
				return token, nil
			}
		}
	}

	// 生成默认 token
	defaultToken := "openclaw-manager-local-token"
	return defaultToken, nil
}

// GetDashboardURL 获取 dashboard URL
func GetDashboardURL() string {
	return fmt.Sprintf("http://localhost:%d", ServicePort)
}
