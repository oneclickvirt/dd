//go:build !windows
// +build !windows

package disk

import (
	"os"
)

// hasRootPermission 检测是否有root权限
func hasRootPermission() bool {
	return os.Getuid() == 0
}
