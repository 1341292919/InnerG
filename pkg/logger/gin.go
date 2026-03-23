package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// GinLogger Gin专用的日志结构体
type GinLogger struct {
	dir        string       // 日志目录
	prefix     string       // 文件名前缀
	currentDay string       // 当前日期（用于判断是否需要切换）
	file       *os.File     // 当前日志文件
	writer     io.Writer    // 输出写入器（文件 + 控制台）
	mu         sync.RWMutex // 读写锁，保证并发安全
}

var GinLog *GinLogger

// InitGinLogger 初始化Gin专用日志系统
// dir: 日志目录，如 "./logs"
// prefix: 文件名前缀，如 "gin"
func InitGinLogger(dir, prefix string) {
	l := &GinLogger{
		dir:    dir,
		prefix: prefix,
	}

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建Gin日志目录失败: %v\n", err)
		panic(err)
	}

	// 打开今天的日志文件
	if err := l.rotate(); err != nil {
		panic(err)
	}

	GinLog = l
}

// Close 关闭当前日志文件
func (l *GinLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// rotate 切换日志文件（内部方法，需要加锁调用）
func (l *GinLogger) rotate() error {
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
func (l *GinLogger) checkRotate() {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != l.currentDay {
		// 日期变了，切换文件
		if err := l.rotate(); err != nil {
			// 切换失败，打印到控制台
			fmt.Fprintf(os.Stderr, "Gin日志切换失败: %v\n", err)
		}
	}
}

// write 写入日志
func (l *GinLogger) write(level, msg string) {
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

// Info 记录信息日志
func (l *GinLogger) Info(v ...interface{}) {
	l.log("INFO", v...)
}

// Error 记录错误日志
func (l *GinLogger) Error(v ...interface{}) {
	l.log("ERROR", v...)
}

// Debug 记录调试日志
func (l *GinLogger) Debug(v ...interface{}) {
	l.log("DEBUG", v...)
}

// Infof 格式化记录信息日志
func (l *GinLogger) Infof(format string, v ...interface{}) {
	l.write("INFO", fmt.Sprintf(format, v...))
}

// Errorf 格式化记录错误日志
func (l *GinLogger) Errorf(format string, v ...interface{}) {
	l.write("ERROR", fmt.Sprintf(format, v...))
}

// Debugf 格式化记录调试日志
func (l *GinLogger) Debugf(format string, v ...interface{}) {
	l.write("DEBUG", fmt.Sprintf(format, v...))
}

// log 多参数拼接日志
func (l *GinLogger) log(level string, v ...interface{}) {
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

// GinWriter 实现 io.Writer 接口，用于 Gin 的日志输出
type GinWriter struct{}

func (gw GinWriter) Write(p []byte) (n int, err error) {
	if GinLog != nil {
		// 去掉末尾的换行符
		msg := strings.TrimSuffix(string(p), "\n")
		if msg != "" {
			GinLog.write("GIN", msg)
		}
	}
	return len(p), nil
}

// GinLoggerMiddleware 返回一个Gin中间件，用于记录详细的请求日志
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录详细日志
		if GinLog != nil {
			GinLog.Infof("[%s] %s | %d | %s | %s",
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				time.Since(start),
				c.ClientIP(),
			)
		}
	}
}

// CloseGinLog 关闭Gin日志文件
func CloseGinLog() {
	if GinLog != nil {
		GinLog.Close()
	}
}
