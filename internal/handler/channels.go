package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iclaw-admin-api/internal/model"
)

// GetChannelsConfig 获取渠道配置 (placeholder)
func GetChannelsConfig(c *gin.Context) {
	c.JSON(http.StatusOK, []model.ChannelConfig{})
}

// SaveChannelConfig 保存渠道配置 (placeholder)
func SaveChannelConfig(c *gin.Context) {
	var req model.ChannelSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Channel config saved (placeholder)"})
}

// ClearChannelConfig 清除渠道配置 (placeholder)
func ClearChannelConfig(c *gin.Context) {
	channelType := c.Param("type")
	c.JSON(http.StatusOK, gin.H{"message": "Channel " + channelType + " config cleared (placeholder)"})
}
