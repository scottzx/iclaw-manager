package service

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"iclaw-admin-api/internal/model"
)

// GetSystemInfo 获取系统信息
func GetSystemInfo() (*model.SystemInfo, error) {
	configDir := GetConfigDir()

	info := &model.SystemInfo{
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		ConfigDir:     configDir,
		OpenClawInstalled: IsOpenClawInstalled(),
	}

	// OS Version
	if runtime.GOOS == "darwin" {
		if out, err := ExecuteCommand("sw_vers", "-productVersion"); err == nil {
			info.OSVersion = strings.TrimSpace(out)
		}
	} else if runtime.GOOS == "linux" {
		if out, err := ExecuteCommand("cat", "/etc/os-release"); err == nil {
			lines := strings.Split(out, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "VERSION=") {
					info.OSVersion = strings.Trim(strings.TrimSpace(line[len("VERSION="):]), "\"")
					break
				}
			}
			if info.OSVersion == "" {
				for _, line := range lines {
					if strings.HasPrefix(line, "PRETTY_NAME=") {
						info.OSVersion = strings.Trim(strings.TrimSpace(line[len("PRETTY_NAME="):]), "\"")
						break
					}
				}
			}
		}
	} else if runtime.GOOS == "windows" {
		if out, err := ExecuteCommand("cmd", "/c", "ver"); err == nil {
			info.OSVersion = strings.TrimSpace(out)
		}
	}

	// OpenClaw version
	if info.OpenClawInstalled {
		if version, err := GetOpenClawVersion(); err == nil {
			info.OpenClawVersion = &version
		}
	}

	// Node version
	if version, err := GetNodeVersion(); err == nil {
		info.NodeVersion = &version
	}

	return info, nil
}

// GetServiceStatus 获取服务状态
func GetServiceStatus() (*model.ServiceStatus, error) {
	pids := GetPidsOnPort(ServicePort)
	running := len(pids) > 0

	status := &model.ServiceStatus{
		Running: running,
		Port:    ServicePort,
	}

	if running && len(pids) > 0 {
		pid := pids[0]
		status.Pid = &pid
	}

	return status, nil
}

// StartService 启动服务
func StartService() (string, error) {
	// 检查是否已经运行
	status, err := GetServiceStatus()
	if err != nil {
		return "", err
	}
	if status.Running {
		return "", fmt.Errorf("服务已在运行中")
	}

	// 检查 openclaw 命令是否存在
	if !IsOpenClawInstalled() {
		return "", fmt.Errorf("找不到 openclaw 命令，请先通过 npm install -g openclaw 安装")
	}

	// 后台启动 gateway
	configDir := GetConfigDir()
	logDir := configDir + "/logs"
	os.MkdirAll(logDir, 0755)

	var cmd string
	if isWindows {
		cmd = fmt.Sprintf(`start /b cmd /c "openclaw gateway > %s\gateway.log 2> %s\gateway.err.log"`, logDir, logDir)
	} else {
		cmd = fmt.Sprintf(`nohup openclaw gateway > %s/gateway.log 2> %s/gateway.err.log &`, logDir, logDir)
	}

	if _, err := ExecuteCommand(cmd); err != nil {
		return "", fmt.Errorf("启动服务失败: %v", err)
	}

	// 等待端口开始监听
	for i := 1; i <= 15; i++ {
		time.Sleep(time.Second)
		if CheckPortInUse(ServicePort) {
			pids := GetPidsOnPort(ServicePort)
			if len(pids) > 0 {
				return fmt.Sprintf("服务已启动，PID: %d", pids[0]), nil
			}
			return "服务已启动", nil
		}
	}

	return "", fmt.Errorf("服务启动超时（15秒），请检查 openclaw 日志")
}

// StopService 停止服务
func StopService() (string, error) {
	pids := GetPidsOnPort(ServicePort)
	if len(pids) == 0 {
		return "服务未在运行", nil
	}

	// 优雅终止
	for _, pid := range pids {
		KillProcess(pid, false)
	}
	time.Sleep(2 * time.Second)

	// 检查是否已停止
	remaining := GetPidsOnPort(ServicePort)
	if len(remaining) == 0 {
		return "服务已停止", nil
	}

	// 强制终止
	for _, pid := range remaining {
		KillProcess(pid, true)
	}
	time.Sleep(time.Second)

	stillRunning := GetPidsOnPort(ServicePort)
	if len(stillRunning) == 0 {
		return "服务已停止", nil
	}

	return "", fmt.Errorf("无法停止服务，仍有进程: %v", stillRunning)
}

// RestartService 重启服务
func RestartService() (string, error) {
	_, _ = StopService()
	time.Sleep(time.Second)
	return StartService()
}

// GetLogs 获取日志
func GetLogs(lines uint32) ([]string, error) {
	if lines == 0 {
		lines = 100
	}

	configDir := GetConfigDir()
	logFiles := []string{
		configDir + "/logs/gateway.log",
		configDir + "/logs/gateway.err.log",
		configDir + "/stderr.log",
		configDir + "/stdout.log",
	}

	var allLines []string
	n := int(lines)

	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			continue
		}

		var cmd string
		if isWindows {
			cmd = fmt.Sprintf("powershell -command \"Get-Content '%s' -Tail %d\"", logFile, n)
		} else {
			cmd = fmt.Sprintf("tail -n %d %s", n, logFile)
		}

		out, err := ExecuteCommand(cmd)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				allLines = append(allLines, trimmed)
			}
		}
	}

	// 去重并保留最后 N 行
	seen := make(map[string]bool)
	var result []string
	for _, line := range allLines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	total := len(result)
	if total > n {
		result = result[total-n:]
	}

	return result, nil
}
