package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// isRoot 判断当前用户是否为 root
func isRoot() bool {
	return os.Geteuid() == 0
}

// CleanDD 删除临时提取出的 coreutils 文件
func CleanDD(tempFile string) error {
	if tempFile == "" {
		return nil // 不需要清理
	}
	// 删除整个临时目录
	return os.RemoveAll(filepath.Dir(tempFile))
}

func main() {
	fmt.Println("")
}
