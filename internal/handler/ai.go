package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iclaw-admin-api/internal/model"
	"iclaw-admin-api/internal/service"
)

// OfficialProviders 官方 Provider 预设
var OfficialProviders = []model.OfficialProvider{
	{
		ID:             "openai",
		Name:           "OpenAI",
		Icon:           "",
		DefaultBaseURL: strPtr("https://api.openai.com"),
		APIType:        "openai-completions",
		SuggestedModels: []model.SuggestedModel{
			{ID: "gpt-4o", Name: "GPT-4o", Recommended: true, ContextWindow: uintPtr(128000), MaxTokens: uintPtr(4096)},
			{ID: "gpt-4o-mini", Name: "GPT-4o mini", Recommended: false, ContextWindow: uintPtr(128000), MaxTokens: uintPtr(16384)},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", Recommended: false, ContextWindow: uintPtr(128000), MaxTokens: uintPtr(4096)},
			{ID: "gpt-3.5-turbo", Name: "GPT-3.5 Turbo", Recommended: false, ContextWindow: uintPtr(16385), MaxTokens: uintPtr(4096)},
		},
		RequiresAPIKey: true,
		DocsURL:        strPtr("https://platform.openai.com/docs"),
	},
	{
		ID:             "anthropic",
		Name:           "Anthropic",
		Icon:           "",
		DefaultBaseURL: strPtr("https://api.anthropic.com"),
		APIType:        "anthropic-messages",
		SuggestedModels: []model.SuggestedModel{
			{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", Recommended: true, ContextWindow: uintPtr(200000), MaxTokens: uintPtr(8192)},
			{ID: "claude-3-5-sonnet-latest", Name: "Claude 3.5 Sonnet", Recommended: false, ContextWindow: uintPtr(200000), MaxTokens: uintPtr(8192)},
			{ID: "claude-3-5-haiku-latest", Name: "Claude 3.5 Haiku", Recommended: false, ContextWindow: uintPtr(200000), MaxTokens: uintPtr(8192)},
			{ID: "claude-3-opus-latest", Name: "Claude 3 Opus", Recommended: false, ContextWindow: uintPtr(200000), MaxTokens: uintPtr(4096)},
		},
		RequiresAPIKey: true,
		DocsURL:        strPtr("https://docs.anthropic.com"),
	},
	{
		ID:             "deepseek",
		Name:           "DeepSeek",
		Icon:           "",
		DefaultBaseURL: strPtr("https://api.deepseek.com"),
		APIType:        "openai-completions",
		SuggestedModels: []model.SuggestedModel{
			{ID: "deepseek-chat", Name: "DeepSeek Chat", Recommended: true, ContextWindow: uintPtr(64000), MaxTokens: uintPtr(4096)},
			{ID: "deepseek-coder", Name: "DeepSeek Coder", Recommended: false, ContextWindow: uintPtr(64000), MaxTokens: uintPtr(4096)},
		},
		RequiresAPIKey: true,
		DocsURL:        strPtr("https://platform.deepseek.com/docs"),
	},
}

func strPtr(s string) *string { return &s }
func uintPtr(u uint32) *uint32 { return &u }

// GetOfficialProviders 获取官方 Provider 列表
func GetOfficialProviders(c *gin.Context) {
	c.JSON(http.StatusOK, OfficialProviders)
}

// GetAIConfig 获取 AI 配置概览
func GetAIConfig(c *gin.Context) {
	overview, err := service.GetAIConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, overview)
}

// SaveProvider 保存 Provider (placeholder)
func SaveProvider(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Provider saved (placeholder)"})
}

// DeleteProvider 删除 Provider (placeholder)
func DeleteProvider(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted (placeholder)"})
}

// SetPrimaryModel 设置主模型 (placeholder)
func SetPrimaryModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Primary model set (placeholder)"})
}

// AddAvailableModel 添加可用模型 (placeholder)
func AddAvailableModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Model added (placeholder)"})
}

// RemoveAvailableModel 移除模型 (placeholder)
func RemoveAvailableModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Model removed (placeholder)"})
}
