package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iclaw-admin-api/internal/model"
)

// GetAgentsList 获取 Agent 列表 (placeholder)
func GetAgentsList(c *gin.Context) {
	c.JSON(http.StatusOK, []model.AgentConfig{})
}

// SaveAgent 保存 Agent (placeholder)
func SaveAgent(c *gin.Context) {
	var req model.AgentSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent saved (placeholder)"})
}

// DeleteAgent 删除 Agent (placeholder)
func DeleteAgent(c *gin.Context) {
	agentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Agent " + agentID + " deleted (placeholder)"})
}

// SetDefaultAgent 设置默认 Agent (placeholder)
func SetDefaultAgent(c *gin.Context) {
	agentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Agent " + agentID + " set as default (placeholder)"})
}
