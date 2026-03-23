package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iclaw-admin-api/internal/model"
)

// GetSkillsList 获取技能列表 (placeholder)
func GetSkillsList(c *gin.Context) {
	c.JSON(http.StatusOK, []model.SkillDefinition{})
}

// InstallSkill 安装技能 (placeholder)
func InstallSkill(c *gin.Context) {
	skillID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Skill " + skillID + " installed (placeholder)"})
}

// UninstallSkill 卸载技能 (placeholder)
func UninstallSkill(c *gin.Context) {
	skillID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Skill " + skillID + " uninstalled (placeholder)"})
}

// SaveSkillConfig 保存技能配置 (placeholder)
func SaveSkillConfig(c *gin.Context) {
	var req model.SkillSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skill config saved (placeholder)"})
}

// InstallCustomSkill 安装自定义技能 (placeholder)
func InstallCustomSkill(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Custom skill installed (placeholder)"})
}
