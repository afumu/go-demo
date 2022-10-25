package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caller",
		MessageKey:     "Msg",
		StacktraceKey:  "St",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   //
		EncodeCaller:   zapcore.FullCallerEncoder,        // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func main() {
	// 动态日志等级
	dynamicLevel := zap.NewAtomicLevel()

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log/agent" + time.Now().Format("2006-01-02") + ".log", // 输出文件
		MaxSize:    100,                                                    // 日志文件最大大小（单位：MB）
		MaxAge:     30,                                                     // 保存日期
		MaxBackups: 3,                                                      // 保留的旧日志文件最大数量
		Compress:   true,                                                   // 压缩归档旧文件
	})

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	core := zapcore.NewTee(
		// 有好的格式、输出控制台、动态等级
		zapcore.NewCore(zapcore.NewConsoleEncoder(NewEncoderConfig()), os.Stdout, dynamicLevel),
		// json格式、输出文件、处定义等级规则
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), w, highPriority),
	)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	sugar := logger.Sugar()

	// 将当前日志等级设置为Debug
	dynamicLevel.SetLevel(zap.InfoLevel)

	sugar.Infof("this is info %v", zap.Int("ID", 1))
	sugar.Errorf("this is info %v", zap.Int("ID", 1))
	sugar.Infof("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	time.Sleep(3 * time.Second)

	sugar.Infof("this is info %v", zap.Int("ID", 1))
	sugar.Errorf("this is info %v", zap.Int("ID", 1))

}
