package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"iam-permission/internal/pkg/config"
	"iam-permission/internal/pkg/constant"
)

var logger *zap.Logger

func InitLogger() error {
	if logger != nil {
		return nil
	}

	logLevel := config.Instance().Common.LogLevel

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		DisableCaller: true,
		Encoding:      "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = l
	return nil
}

func fromCtx(ctx context.Context) *zap.Logger {
	l := ctx.Value(constant.CtxLoggerKey)
	if l == nil {
		l = ctx.Value(ctxLoggerKey(constant.CtxLoggerKey))
	}
	if l == nil {
		return logger
	}
	zlogger, ok := l.(*zap.Logger)
	if !ok {
		return logger
	}
	return zlogger
}

func Error(ctx context.Context, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Errorf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Warnf(format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	l := fromCtx(ctx)
	l.Sugar().Debugf(format, args...)
}

func L() *zap.Logger {
	return logger
}
