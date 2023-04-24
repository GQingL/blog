// Copyright 2022 Innkeeper Belm(梁广庆) &lt;138521257@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/guangqingliang/blog

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

// 确保 zapLogger 实现了 Logger 接口
var _ Logger = &zapLogger{}

// zapLogger 是 logger 接口的具体实现，它底层封装了zap。logger
type zapLogger struct {
	z *zap.Logger
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Infow(msg, keysAndValues...)
}

func (z zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.z.Sugar().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Warnw(msg, keysAndValues...)
}

func (z zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.z.Sugar().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Errorw(msg, keysAndValues...)
}

func (z zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.z.Sugar().Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Panicw(msg, keysAndValues...)
}

func (z zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.z.Sugar().Panicw(msg, keysAndValues)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Fatalw(msg, keysAndValues...)
}

func (z zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.z.Sugar().Fatalw(msg, keysAndValues)
}

func Sync() {
	std.Sync()
}

func (z zapLogger) Sync() {
	_ = z.z.Sync()
}

var (
	mu sync.Mutex

	// std 定义了默认的全局Logger
	std = NewLogger(NewOptions())
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger(opts)
}

// NewLogger 根据传入的opts 创建 Logger
func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将文本格式的日志级别，例如 info 转换为zapcore.Level 类型以供后面使用
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果指定了非法的日志级别，则默认使用info级别
		zapLevel = zapcore.InfoLevel
	}

	// 创建一个默认的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义 MessageKey 为 message, timestamp 语意更明确
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"

	// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// 指定time.Duration序列化函数，将time.Duration序列化为经过的毫秒数的浮点数
	// 毫秒数比默认的秒数更精确
	encoderConfig.EncodeDuration = func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendFloat64(float64(duration) / float64(time.Millisecond))
	}

	cfg := &zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
		// 是否再日志中显示调用日志所在的文件和行号
		DisableCaller: opts.DisableCaller,
		// 是否禁止再panic及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStacktrace,
		// 指定日志显示格式，可选值：console，json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 限制日志输出位置
		OutputPaths: opts.OutputPaths,
		// 设置zap内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}
	// 使用cfg创建*zap.Logger对象
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if nil != err {
		panic(err)
	}

	logger := &zapLogger{z: z}
	zap.RedirectStdLog(z)
	return logger
}
