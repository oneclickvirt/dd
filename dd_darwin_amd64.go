//go:build darwin && amd64
// +build darwin,amd64

package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed bin/coreutils-darwin-amd64
var binFiles embed.FS

// GetDD 获取与当前系统匹配的 dd 二进制文件并返回路径
func GetDD() (string, error) {
	binaryName := "coreutils-darwin-amd64"

	// 检查系统是否有原生 dd 命令
	if _, err := exec.LookPath("dd"); err == nil {
		return "dd", nil // 返回系统原生命令
	}

	// 创建临时目录存放二进制文件
	tempDir, err := os.MkdirTemp("", "ddwrapper")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 读取嵌入的二进制文件
	binPath := filepath.Join("bin", binaryName)
	fileContent, err := binFiles.ReadFile(binPath)
	if err != nil {
		return "", fmt.Errorf("读取嵌入的 coreutils 二进制文件失败: %v", err)
	}

	// 写入临时文件
	tempFile := filepath.Join(tempDir, binaryName)
	if err := os.WriteFile(tempFile, fileContent, 0755); err != nil {
		return "", fmt.Errorf("写入临时文件失败: %v", err)
	}

	return tempFile, nil
}

// ExecuteDD 执行 dd 命令
func ExecuteDD(ddPath string, args []string) error {
	var cmd *exec.Cmd
	if ddPath == "dd" {
		// 使用系统 dd
		cmd = exec.Command(ddPath, args...)
	} else {
		// 使用提取的 coreutils dd
		ddCmd := fmt.Sprintf("%s dd %s", ddPath, strings.Join(args, " "))
		cmd = exec.Command("sh", "-c", ddCmd)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
