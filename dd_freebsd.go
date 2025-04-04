//go:build freebsd
// +build freebsd

package dd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// //go:embed bin/coreutils-linux-amd64 bin/coreutils-linux-musl-amd64
// var binFiles embed.FS

// GetDD 返回适用于当前系统的 dd 命令字符串（带不带 sudo）和临时文件路径（用于清理）
func GetDD() (ddCmd string, tempFile string, err error) {
	// 优先尝试 sudo dd 是否可用
	if path, err := exec.LookPath("dd"); err == nil {
		testCmd := exec.Command("sudo", path, "--help")
		if err := testCmd.Run(); err == nil {
			return "sudo dd", "", nil
		}
		// 如果 sudo dd 不可用，则尝试直接使用 dd
		testCmd = exec.Command(path, "--help")
		if err := testCmd.Run(); err == nil {
			return "dd", "", nil
		}
	}
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ddwrapper")
	if err != nil {
		return "", "", fmt.Errorf("创建临时目录失败: %v", err)
	}
	// // 尝试使用 glibc 版本
	// binName := "coreutils-linux-amd64"
	// binPath := filepath.Join("bin", binName)
	// fileContent, err := binFiles.ReadFile(binPath)
	// if err == nil {
	// 	tempFile = filepath.Join(tempDir, binName)
	// 	if err := os.WriteFile(tempFile, fileContent, 0755); err == nil {
	// 		// 先尝试 sudo 运行
	// 		testCmd := exec.Command("sudo", tempFile, "--version")
	// 		if err := testCmd.Run(); err == nil {
	// 			return fmt.Sprintf("sudo %s dd", tempFile), tempFile, nil
	// 		}
	// 		// 如果 sudo 运行失败，尝试直接运行
	// 		testCmd = exec.Command(tempFile, "--version")
	// 		if err := testCmd.Run(); err == nil {
	// 			return fmt.Sprintf("%s dd", tempFile), tempFile, nil
	// 		}
	// 	}
	// }
	// // 尝试使用 musl 版本
	// binName = "coreutils-linux-musl-amd64"
	// binPath = filepath.Join("bin", binName)
	// fileContent, err = binFiles.ReadFile(binPath)
	// if err != nil {
	// 	return "", "", fmt.Errorf("读取嵌入的 coreutils 二进制文件失败: %v", err)
	// }
	// tempFile = filepath.Join(tempDir, binName)
	// if err := os.WriteFile(tempFile, fileContent, 0755); err != nil {
	// 	return "", "", fmt.Errorf("写入临时文件失败: %v", err)
	// }
	// // 先尝试 sudo 运行
	// testCmd := exec.Command("sudo", tempFile, "--version")
	// if err := testCmd.Run(); err == nil {
	// 	return fmt.Sprintf("sudo %s dd", tempFile), tempFile, nil
	// }
	// // 如果 sudo 运行失败，尝试直接运行
	// testCmd = exec.Command(tempFile, "--version")
	// if err := testCmd.Run(); err == nil {
	// 	return fmt.Sprintf("%s dd", tempFile), tempFile, nil
	// }
	return "", "", fmt.Errorf("无法找到可用的 dd 命令")
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
