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

// GetDD 获取与当前系统匹配的 dd 命令（包含是否需要 sudo），并返回完整执行路径或命令
func GetDD() (string, string, error) {
	binaryName := "coreutils-darwin-amd64"
	useSudo := !isRoot()
	// 检查系统是否有原生 dd 命令
	if _, err := exec.LookPath("dd"); err == nil {
		if useSudo {
			return "sudo dd", "", nil
		}
		return "dd", "", nil
	}
	// 创建临时目录存放二进制文件
	tempDir, err := os.MkdirTemp("", "ddwrapper")
	if err != nil {
		return "", "", fmt.Errorf("创建临时目录失败: %v", err)
	}
	// 读取嵌入的二进制文件
	binPath := filepath.Join("bin", binaryName)
	fileContent, err := binFiles.ReadFile(binPath)
	if err != nil {
		return "", "", fmt.Errorf("读取嵌入的 coreutils 二进制文件失败: %v", err)
	}
	// 写入临时文件
	tempFile := filepath.Join(tempDir, binaryName)
	if err := os.WriteFile(tempFile, fileContent, 0755); err != nil {
		return "", "", fmt.Errorf("写入临时文件失败: %v", err)
	}
	// 返回完整的命令
	if useSudo {
		return fmt.Sprintf("sudo %s dd", tempFile), tempFile, nil
	}
	return fmt.Sprintf("%s dd", tempFile), tempFile, nil
}

// ExecuteDD 执行拼好的 dd 命令字符串（包括 sudo、dd 等）
func ExecuteDD(ddCmd string, args []string) error {
	// 拼接命令字符串
	fullCmd := fmt.Sprintf("%s %s", ddCmd, strings.Join(args, " "))
	cmd := exec.Command("sh", "-c", fullCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
