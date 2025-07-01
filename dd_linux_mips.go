//go:build linux && mips
// +build linux,mips

package dd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed bin/dd-linux-mips
var binFiles embed.FS

// GetDD 返回适用于当前系统的 dd 命令字符串（带不带 sudo）和临时文件路径（用于清理）
func GetDD() (ddCmd string, tempFile string, err error) {
	var errors []string
	// 1. 尝试系统自带 dd
	if path, lookErr := exec.LookPath("dd"); lookErr == nil {
		if !hasRootPermission() {
			// 确保 testCmd 被初始化
			testCmd := exec.Command("sudo", path, "--help")
			if runErr := testCmd.Run(); runErr == nil {
				return "sudo dd", "", nil
			} else {
				errors = append(errors, fmt.Sprintf("sudo dd 测试失败: %v", runErr))
			}
		}
		// 直接尝试 dd
		testCmd := exec.Command(path, "--help")
		if runErr := testCmd.Run(); runErr == nil {
			return "dd", "", nil
		} else {
			errors = append(errors, fmt.Sprintf("dd 直接运行失败: %v", runErr))
		}
	} else {
		errors = append(errors, fmt.Sprintf("无法找到 dd: %v", lookErr))
	}
	// 2. 创建临时目录
	tempDir, tempErr := os.MkdirTemp("", "ddwrapper")
	if tempErr != nil {
		return "", "", fmt.Errorf("创建临时目录失败: %v", tempErr)
	}
	// 3. 尝试使用 glibc 版本 coreutils
	binName := "dd-linux-mips"
	binPath := filepath.Join("bin", binName)
	fileContent, readErr := binFiles.ReadFile(binPath)
	if readErr == nil {
		tempFile = filepath.Join(tempDir, binName)
		writeErr := os.WriteFile(tempFile, fileContent, 0755)
		if writeErr == nil {
			if !hasRootPermission() {
				// 确保 testCmd 被初始化
				testCmd := exec.Command("sudo", tempFile, "--version")
				if runErr := testCmd.Run(); runErr == nil {
					return fmt.Sprintf("sudo %s", tempFile), tempFile, nil
				} else {
					errors = append(errors, fmt.Sprintf("sudo %s 运行失败: %v", tempFile, runErr))
				}
			}
			// 直接尝试
			testCmd := exec.Command(tempFile, "--version")
			if runErr := testCmd.Run(); runErr == nil {
				return tempFile, tempFile, nil
			} else {
				errors = append(errors, fmt.Sprintf("%s 运行失败: %v", tempFile, runErr))
			}
		} else {
			errors = append(errors, fmt.Sprintf("写入临时文件失败 (%s): %v", tempFile, writeErr))
		}
	} else {
		errors = append(errors, fmt.Sprintf("读取嵌入的 coreutils glibc 版本失败: %v", readErr))
	}
	// 返回所有错误信息
	return "", "", fmt.Errorf("无法找到可用的 dd 命令:\n%s", strings.Join(errors, "\n"))
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
