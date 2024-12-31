package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateConfig struct {
	MaxSizeMB   int
	MaxAgeDays  int
	MaxBackups  int
	GzipArchive bool
}

var defaultLogger = zap.Must(zap.NewDevelopment())

func NewLogger(debug bool, logFile string, enableRotation bool) *zap.Logger {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: debug,
		OutputPaths: []string{},
	}

	if debug {
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Encoding = "console"
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
	} else {
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.Encoding = "json"
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}

	return zap.Must(cfg.Build(zap.AddStacktrace(zap.DPanicLevel)))
}

func NewLoggerWithLogRotate(debug bool, logFile string, rotateCfg RotateConfig) *zap.Logger {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: debug,
		OutputPaths: []string{},
	}

	if debug {
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    rotateCfg.MaxSizeMB,
		MaxAge:     rotateCfg.MaxAgeDays,
		Compress:   rotateCfg.GzipArchive,
		MaxBackups: rotateCfg.MaxBackups,
	})
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg.EncoderConfig), ws, zap.InfoLevel)
	return zap.New(core, zap.AddStacktrace(zap.DPanicLevel))
}

func Debug(msg string, fields ...zapcore.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	defaultLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	defaultLogger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	defaultLogger.Panic(msg, fields...)
}
