package extensions

import (
	"github.com/goravel/framework/contracts/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
}

func (z *ZapLogger) Handle(channel string) (log.Hook, error) {
	return NewZapHook(channel), nil
}

func NewZapHook(channel string) log.Hook {
	zp, _ := zap.NewProduction()
	return &ZapHook{
		channel,
		zp,
	}
}

type ZapHook struct {
	channel string
	zp      *zap.Logger
}

func (z *ZapHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarningLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

func (z *ZapHook) Fire(entry log.Entry) error {
	ch := z.zp.Named(z.channel).WithOptions(zap.AddCallerSkip(10), zap.AddStacktrace(zapcore.ErrorLevel))
	options := []zap.Field{
		zap.Any("data", entry.Data()),
	}
	switch entry.Level() {
	default:
		ch.Debug(entry.Message(), options...)
	case log.InfoLevel:
		ch.Info(entry.Message(), options...)
	case log.WarningLevel:
		ch.Warn(entry.Message(), options...)
	case log.ErrorLevel:
		ch.Error(entry.Message(), options...)
	case log.FatalLevel:
		ch.Fatal(entry.Message(), options...)
	case log.PanicLevel:
		ch.Panic(entry.Message(), options...)
	}
	defer ch.Sync()
	return nil
}
