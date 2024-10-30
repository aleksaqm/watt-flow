package util

import (
	"context"
	"fmt"
	"time"
	"watt-flow/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	*zap.SugaredLogger
}

type GinLogger struct {
	*Logger
}

type GormLogger struct {
	*Logger
	gormlogger.Config
}

var (
	globalLogger *Logger
	zapLogger    *zap.Logger
)

func (l *Logger) GetGinLogger() GinLogger {
	logger := zapLogger.WithOptions(zap.WithCaller(false))
	return GinLogger{
		Logger: newSugaredLogger(logger),
	}
}

func (l Logger) GetGormLogger() *GormLogger {
	logger := zapLogger.WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)

	return &GormLogger{
		Logger: newSugaredLogger(logger),
		Config: gormlogger.Config{
			LogLevel: gormlogger.Info,
		},
	}
}

func newSugaredLogger(logger *zap.Logger) *Logger {
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

func NewLogger(env *config.Environment) Logger {
	config := zap.NewDevelopmentConfig()
	logOutput := env.LogOutput

	if env.Environment == "development" {
		fmt.Println("encode level")
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if env.Environment == "production" && logOutput != "" {
		config.OutputPaths = []string{logOutput}
	}

	logLevel := env.LogLevel
	level := zap.PanicLevel
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zap.PanicLevel
	}
	config.Level.SetLevel(level)

	zapLogger, _ = config.Build()
	logger := newSugaredLogger(zapLogger)

	return *logger
}

func (l GinLogger) Write(p []byte) (n int, err error) {
	l.Info(string(p))
	return len(p), nil
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info prints info
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Debugf(str, args...)
	}
}

// Warn prints warn messages
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(str, args...)
	}
}

// Error prints error messages
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(str, args...)
	}
}

// Trace prints trace messages
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	if l.LogLevel >= gormlogger.Info {
		sql, rows := fc()
		l.Debug("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}

	if l.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.SugaredLogger.Warn("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}

	if l.LogLevel >= gormlogger.Error {
		sql, rows := fc()
		l.SugaredLogger.Error("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}
}
