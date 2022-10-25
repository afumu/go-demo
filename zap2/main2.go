package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

var logger *zap.Logger

// logpath 日志文件路径
// loglevel 日志级别
func InitLogger(loglevel string) {
	// 日志分割
	hook := lumberjack.Logger{
		Filename: "log/agent" + time.Now().Format("2006-01-02") + ".log", // 输出文件
		MaxSize:  10,                                                     // 每个日志文件保存10M，默认 100M
		MaxAge:   7,                                                      // 保留7天，默认不限
		Compress: true,                                                   // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger = zap.New(core, caller, development, filed)
	logger.Info("DefaultLogger init success")
}
func main() {
	// 历史记录日志名字为：all.log，服务重新启动，日志会追加，不会删除
	InitLogger("debug")
	// 强结构形式
	logger.Info("这是一个测试log")
	// 必须 key-value 结构形式 性能下降一点
	logger.Sugar().Infow("test-",
		"string", "string",
		"int", 1,
		"time", time.Second,
	)
}
