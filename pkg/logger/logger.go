package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	dir        string       // 日志目录
	prefix     string       // 文件名前缀
	currentDay string       // 当前日期（用于判断是否需要切换）
	file       *os.File     // 当前日志文件
	writer     io.Writer    // 输出写入器（文件 + 控制台）
	mu         sync.RWMutex // 读写锁，保证并发安全
}

var Log *Logger

// InitLogger 初始化日志系统
// dir: 日志目录，如 "./logs"
// prefix: 文件名前缀，如 "app"
func InitLogger(dir, prefix string) {
	l := &Logger{
		dir:    dir,
		prefix: prefix,
	}

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("创建日志目录失败: %w", err)
		panic(err)
	}

	// 打开今天的日志文件
	if err := l.rotate(); err != nil {
		panic(err)
	}

	Log = l
	return
}

// Close 关闭当前日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// 切换日志文件（内部方法，需要加锁调用）
func (l *Logger) rotate() error {
	today := time.Now().Format("2006-01-02")
	filename := filepath.Join(l.dir, fmt.Sprintf("%s-%s.log", l.prefix, today))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	// 关闭旧文件（如果存在）
	if l.file != nil {
		l.file.Close()
	}

	l.currentDay = today
	l.file = file
	l.writer = io.MultiWriter(file, os.Stdout)

	return nil
}

// checkRotate 检查是否需要切换文件（每天第一次写入时触发）
func (l *Logger) checkRotate() {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != l.currentDay {
		// 日期变了，切换文件
		if err := l.rotate(); err != nil {
			// 切换失败，打印到控制台
			fmt.Fprintf(os.Stderr, "日志切换失败: %v\n", err)
		}
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.log("INFO", v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log("ERROR", v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.log("DEBUG", v...)
}

// ============ 格式化输出版本（新增） ============

func (l *Logger) Infof(format string, v ...interface{}) {
	l.write("INFO", fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write("ERROR", fmt.Sprintf(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write("DEBUG", fmt.Sprintf(format, v...))
}

// ============ 内部方法 ============

// 多参数拼接
func (l *Logger) log(level string, v ...interface{}) {
	l.mu.RLock()
	today := time.Now().Format("2006-01-02")
	if today != l.currentDay {
		l.mu.RUnlock()
		l.checkRotate()
		l.mu.RLock()
	}
	defer l.mu.RUnlock()

	var parts []string
	for _, item := range v {
		parts = append(parts, fmt.Sprint(item))
	}
	msg := strings.Join(parts, "")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, msg)
	l.writer.Write([]byte(line))
}

// 格式化输出
func (l *Logger) write(level, msg string) {
	l.mu.RLock()
	today := time.Now().Format("2006-01-02")
	if today != l.currentDay {
		l.mu.RUnlock()
		l.checkRotate()
		l.mu.RLock()
	}
	defer l.mu.RUnlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, msg)
	l.writer.Write([]byte(line))
}

func CloseAll() {
	if Log != nil {
		Log.Close()
	}
	CloseGinLog() // 调用 gin_logger.go 中的关闭函数
}
