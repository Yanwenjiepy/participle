package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.Logger

func InitLog(logPath, logLevel string) *zap.Logger {

	hook := lumberjack.Logger{
		Filename: logPath,
		MaxSize:  500,
		MaxAge:   30,
		Compress: true,
	}

	w := zapcore.AddSync(&hook)

	var level zapcore.Level

	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	consoleError := zapcore.Lock(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	fileCore := zapcore.NewCore(fileEncoder, w, level)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleError, level)
	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core, zap.AddCaller())
	return logger

}
