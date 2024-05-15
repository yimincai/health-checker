package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CronLogger struct {
	Logger *zap.Logger
}

var (
	loggerOnce sync.Once
	zapLogger  *zap.Logger
)

func New() {
	ljl := &lumberjack.Logger{
		Filename:   "./bot_files/healthchecker.log",
		MaxSize:    10, // megabytes
		MaxAge:     30, // days
		MaxBackups: 3,
		LocalTime:  true,
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

func NewCronLogger() *CronLogger {
	return &CronLogger{zapLogger}
}

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	fields := make([]zap.Field, 0, len(keysAndValues))
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
		}
	}
	zapLogger.Info(msg, fields...)
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := make([]zap.Field, 0, len(keysAndValues)+1)
	fields = append(fields, zap.Error(err))
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
		}
	}
	zapLogger.Error(msg, fields...)
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
