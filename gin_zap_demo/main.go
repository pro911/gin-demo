package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/pro911/gin-demo/gin_zap_demo/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var logger *zap.Logger
var loggers *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriterAll() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./debug.log", //指定写入日志的文件路径
		MaxSize:    1,             //最大/M
		MaxAge:     5,             //备份数量 多少个文件
		MaxBackups: 2,             //最大备份天数
		LocalTime:  false,
		Compress:   false, //是否压缩
	}
	// 利用io.MultiWriter支持文件和终端两个输出目标
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

func getLogWriterErr() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./err.log", //指定写入日志的文件路径
		MaxSize:    1,           //最大/M
		MaxAge:     5,           //备份数量 多少个文件
		MaxBackups: 2,           //最大备份天数
		LocalTime:  false,
		Compress:   false, //是否压缩
	}
	// 利用io.MultiWriter支持文件和终端两个输出目标
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

func InitLogger() {
	encoder := getEncoder()

	// test.log记录全量日志 & 错误日志
	writeSyncerAll := getLogWriterAll()
	writeSyncerErr := getLogWriterErr()

	// test.err.log记录ERROR级别的日志
	coreAll := zapcore.NewCore(encoder, writeSyncerAll, zapcore.DebugLevel)
	coreErr := zapcore.NewCore(encoder, writeSyncerErr, zapcore.ErrorLevel)
	//logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	//使用NewTee将多个文件合并到core
	core := zapcore.NewTee(coreAll, coreErr)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	loggers = logger.Sugar()
}

func main() {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello liwenzhou.com!")
	})
	r.Run(":8000")
}
