package dd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CleanDD 删除临时提取出的 coreutils 文件
func CleanDD(tempFile string) error {
	if tempFile == "" {
		return nil // 不需要清理
	}
	// 删除整个临时目录
	return os.RemoveAll(filepath.Dir(tempFile))
}