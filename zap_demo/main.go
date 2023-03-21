package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
)

var logger *zap.Logger
var loggers *zap.SugaredLogger

//func InitLogger() {
//	logger, _ = zap.NewProduction()
//	loggers = logger.Sugar()
//}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewJSONEncoder(encoderConfig)

	return zapcore.NewConsoleEncoder(encoderConfig)
}

//func getLogWriterAll() zapcore.WriteSyncer {
//	file, _ := os.OpenFile("./debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//	// 利用io.MultiWriter支持文件和终端两个输出目标
//	ws := io.MultiWriter(file, os.Stdout)
//	return zapcore.AddSync(ws)
//}

//func getLogWriterErr() zapcore.WriteSyncer {
//	file, _ := os.OpenFile("./err.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//	// 利用io.MultiWriter支持文件和终端两个输出目标
//	ws := io.MultiWriter(file, os.Stdout)
//	return zapcore.AddSync(ws)
//}

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
		MaxSize:    2,           //最大
		MaxAge:     5,
		MaxBackups: 5,
		LocalTime:  false,
		Compress:   false,
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

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching url...", zap.String("url", url), zap.Error(err))
		loggers.Error("Error fetching url...", zap.String("url", url), zap.Error(err))
	} else {
		logger.Info("Success...", zap.String("statusCode", resp.Status), zap.String("url", url))
		loggers.Info("Success...", zap.String("statusCode", resp.Status), zap.String("url", url))
		defer resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer logger.Sync()
	defer loggers.Sync()
	for i := 0; i < 100000; i++ {
		simpleHttpGet("https://www.baidu.com")
	}
}
