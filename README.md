# dd

一个嵌入dd依赖的golang库(A golang library with embedded dd dependencies)

## 关于 About

这个库提供了对dd命令行工具的Go语言封装。dd是一个广泛使用的UNIX命令行实用程序，用于进行底层数据复制和转换操作，常用于磁盘克隆、数据备份、设备格式化等操作。

This library provides a Go wrapper for the dd command-line tool. dd is a widely used UNIX command-line utility for low-level data copying and conversion operations, commonly used for disk cloning, data backup, device formatting, and other operations.

## 特性 Features

- 支持多平台：Linux (amd64, arm64, 386, arm等), macOS (amd64, arm64), Windows (amd64, 386)
- 自动检测并使用系统安装的dd，或使用嵌入的二进制文件
- 自动权限检测和sudo处理
- 自动清理临时文件
- Multi-platform support: Linux (amd64, arm64, 386, arm, etc.), macOS (amd64, arm64), Windows (amd64, 386)
- Automatically detects and uses system-installed dd, or uses embedded binaries
- Automatic permission detection and sudo handling
- Automatic cleanup of temporary files

## 安装 Installation

```bash
go get github.com/oneclickvirt/dd@v0.0.2-20250808062818
```

## 使用方法 Usage

```go
package main

import (
    "log"
    "github.com/oneclickvirt/dd"
)

func main() {
    // 获取dd命令路径
    // Get dd command path
    ddCmd, tempFile, err := dd.GetDD()
    if err != nil {
        log.Fatalf("Failed to get dd: %v", err)
    }

    // 如果使用了临时文件，确保清理
    // Clean up temporary files if used
    if tempFile != "" {
        defer dd.CleanDD(tempFile)
    }

    // 执行dd命令（例如：创建一个10MB的零文件）
    // Execute dd command (example: create a 10MB zero file)
    err = dd.ExecuteDD(ddCmd, []string{"if=/dev/zero", "of=test.img", "bs=1M", "count=10"})
    if err != nil {
        log.Fatalf("Failed to execute dd: %v", err)
    }
}
```

## API文档 API Documentation

### 函数 Functions

#### `GetDD() (string, string, error)`

获取可用的dd命令路径。

Get the available dd command path.

**返回值 Returns:**
- `string`: dd命令的路径（可能包含sudo前缀）(Path to the dd command, may include sudo prefix)
- `string`: 临时文件路径（如果使用了嵌入的二进制文件）(Temporary file path if using embedded binary)
- `error`: 错误信息 (Error information)

**说明 Description:**
- 优先尝试使用系统已安装的dd命令 (Prioritizes system-installed dd command)
- 如果系统dd不可用，则提取并使用内嵌的dd二进制文件 (If system dd is unavailable, extracts and uses embedded dd binary)
- 自动检测root权限，必要时添加sudo前缀 (Automatically detects root permissions and adds sudo prefix when necessary)

#### `ExecuteDD(ddPath string, args []string) error`

执行dd命令。

Execute the dd command.

**参数 Parameters:**
- `ddPath`: dd命令的路径 (Path to the dd command)
- `args`: 传递给dd的参数 (Arguments to pass to dd)

**返回值 Returns:**
- `error`: 错误信息 (Error information)

**示例 Examples:**
```go
// 复制磁盘分区
// Copy disk partition
dd.ExecuteDD(ddCmd, []string{"if=/dev/sda1", "of=/dev/sdb1", "bs=4M"})

// 创建交换文件
// Create swap file
dd.ExecuteDD(ddCmd, []string{"if=/dev/zero", "of=/swapfile", "bs=1G", "count=2"})

// 备份MBR
// Backup MBR
dd.ExecuteDD(ddCmd, []string{"if=/dev/sda", "of=mbr_backup.img", "bs=512", "count=1"})
```

#### `CleanDD(tempFile string) error`

清理临时文件。

Clean up temporary files.

**参数 Parameters:**
- `tempFile`: 需要清理的临时文件路径 (Path to temporary file to clean up)

**返回值 Returns:**
- `error`: 错误信息 (Error information)

**说明 Description:**
- 如果tempFile为空字符串，则无需清理 (No cleanup needed if tempFile is empty string)
- 删除整个临时目录及其内容 (Removes entire temporary directory and its contents)

## 平台支持 Platform Support

库包含了以下平台的预编译二进制文件：

The library includes precompiled binaries for the following platforms:

### Linux
- amd64 (x86_64)
- 386 (x86 32-bit)
- arm64 (ARMv8)
- arm (ARMv7)
- riscv64 (RISC-V 64-bit)
- ppc64le (PowerPC64 little-endian)
- ppc64 (PowerPC64 big-endian)
- mips64le (MIPS64 little-endian)
- mips64 (MIPS64 big-endian)
- mipsle (MIPS little-endian)
- mips (MIPS big-endian)

### macOS
- amd64 (Intel)
- arm64 (Apple Silicon)

### Windows
- amd64 (x86_64)
- 386 (x86 32-bit)

### 其他系统 Other Systems
- FreeBSD
- OpenBSD
- 其他类Unix系统（回退到系统dd）(Other Unix-like systems, fallback to system dd)

## 二进制文件来源 Binary Sources

本库内嵌的dd二进制文件来源于：

The embedded dd binaries in this library are sourced from:

- **原版 Original**: [GNU Coreutils](https://github.com/coreutils/coreutils)
- **Rust重构版本 Rust Rewrite**: [uutils/coreutils](https://github.com/uutils/coreutils)

## 许可证 License

MIT License