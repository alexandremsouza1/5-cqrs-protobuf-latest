package logger

import (
	"context"
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var logger *zap.Logger

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return logger.Sugar()
	}
	if ctxlogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxlogger.Sugar()
	} else {
		return logger.Sugar()
	}
}

func GetSugarLogger(ctx context.Context) *zap.Logger {
	return logger
}

func Info(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Warnf(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Warn(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Errorf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Error(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Errorf(format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Fatal(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Debugf(format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Debug(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	stampObservabilityFields(ctx).Panicf(format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	stampObservabilityFields(ctx).Panic(args...)
}

func stampObservabilityFields(ctx context.Context) *zap.SugaredLogger {
	ipAddress := getIpAddressOrigin(ctx)
	if ipAddress != "" {
		return WithContext(ctx).With(
			zap.Field{Key: "ip-origin", Type: zapcore.StringType, String: ipAddress},
		)
	}
	return WithContext(ctx)
}

func getIpAddressOrigin(ctx context.Context) string {
	var ipAddress string
	if ctx != nil {
		peer, ok := peer.FromContext(ctx)
		if ok {
			ipAddress = peer.Addr.String()
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				realIp := md.Get("X-Real-IP")
				if realIp != nil {
					ipAddress = strings.Join(realIp, "")
				}
			}
		}
	}
	return ipAddress
}

func InitPlain() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Print("Unable to initialize logging subsystem")
	}
}

func InitZapCores(enc zapcore.Encoder, enab zapcore.LevelEnabler, wSyncers ...zapcore.WriteSyncer) {
	cores := make([]zapcore.Core, 0)
	for _, ws := range wSyncers {
		cores = append(cores, zapcore.NewCore(enc, ws, enab))
	}
	logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
}
