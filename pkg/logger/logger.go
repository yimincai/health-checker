package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Logger *zap.Logger
}

var loggerOnce sync.Once
var zapLogger *zap.Logger

func init() {
	ljl := &lumberjack.Logger{
		Filename:   "./bot_files/healthchecker.log",
		MaxSize:    10, // megabytes
		MaxAge:     30, // days
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}

	writeSyncer := zapcore.AddSync(ljl)

	loggerOnce.Do(func() {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(cfg)
		core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
		zapLogger = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller())
		zap.ReplaceGlobals(zapLogger)
	})
}

func GetInstance() *Logger {
	return &Logger{zapLogger}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, keysAndValues...))
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.Logger.Error(fmt.Sprintf(fmt.Sprintf("%s error: %s", msg, err), keysAndValues...))
}

func Sugar() *zap.SugaredLogger {
	return zapLogger.Sugar()
}

func Infof(template string, args ...interface{}) {
	zapLogger.Info(fmt.Sprintf(template, args...))
}

func Info(args ...interface{}) {
	zapLogger.Info(fmt.Sprintf("%v", args...))
}

func Debugf(template string, args ...interface{}) {
	zapLogger.Debug(fmt.Sprintf(template, args...))
}

func Debug(args ...interface{}) {
	zapLogger.Debug(fmt.Sprintf("%v", args...))
}

func Warnf(template string, args ...interface{}) {
	zapLogger.Warn(fmt.Sprintf(template, args...))
}

func Warn(args ...interface{}) {
	zapLogger.Warn(fmt.Sprintf("%v", args...))
}

func Errorf(template string, args ...interface{}) {
	zapLogger.Error(fmt.Sprintf(template, args...))
}

func Error(args ...interface{}) {
	zapLogger.Error(fmt.Sprintf("%v", args...))
}

func Fatalf(template string, args ...interface{}) {
	zapLogger.Fatal(fmt.Sprintf(template, args...))
}

func Fatal(args ...interface{}) {
	zapLogger.Fatal(fmt.Sprintf("%v", args...))
}

func Panicf(template string, args ...interface{}) {
	zapLogger.Panic(fmt.Sprintf(template, args...))
}

func Panic(args ...interface{}) {
	zapLogger.Panic(fmt.Sprintf("%v", args...))
}
