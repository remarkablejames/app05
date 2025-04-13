package logger

import (
	"app05/internal/core/application/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"

	// Background colors
	bgBlack   = "\033[40m"
	bgRed     = "\033[41;1m" // Bright red background
	bgGreen   = "\033[42;1m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
	bgWhite   = "\033[47m"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() contracts.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "console" // ⚡ Enable human-readable console output

	// Custom encoder config
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    CustomLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, _ := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	return &ZapLogger{logger}
}

// CustomLevelEncoder adds colors to log levels
func CustomLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var level string
	switch l {
	case zapcore.DebugLevel:
		level = bgBlue + colorWhite + " DEBUG " + colorReset
	case zapcore.InfoLevel:
		level = bgGreen + colorBlack + " INFO  " + colorReset
	case zapcore.WarnLevel:
		level = bgYellow + colorBlack + " WARN  " + colorReset
	case zapcore.ErrorLevel:
		level = bgRed + colorBlack + " ERROR " + colorReset
	case zapcore.PanicLevel:
		level = bgMagenta + colorWhite + " PANIC " + colorReset
	case zapcore.FatalLevel:
		level = bgRed + colorWhite + " FATAL " + colorReset
	default:
		level = bgBlack + colorWhite + "UNKNOWN" + colorReset
	}
	enc.AppendString(level)
}

// Updated logger methods to match the interface
func (z *ZapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Debug(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...interface{}) {
	z.logger.Info(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...interface{}) {
	z.logger.Warn(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...interface{}) {
	z.logger.Error(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Panic(msg string, fields ...interface{}) {
	z.logger.Panic(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...interface{}) {
	z.logger.Fatal(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

// Helper function to convert interface{} to zap.Field
func convertToZapFields(fields []interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}
			zapFields = append(zapFields, zap.Any(key, fields[i+1]))
		}
	}
	return zapFields
}
