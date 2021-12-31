package logger

import (
	"errors"
	"github.com/xfyun/lumberjack-ccr"
	"github.com/xfyun/xns/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var globalpid int64
var async = true
var cacheMaxCount = -1
var batchSize = 16 * 1024
var defaultWash = 60

func init() {
	globalpid = int64(os.Getpid())
}

type Logger struct {
	lumberjack *lumberjack.Logger
	logger     *zap.SugaredLogger
}


type LogConf = conf.LogConf


func newLocalLog(lc *conf.LogConf) (*Logger, error) {
	var localLog   *Logger
	//	fmt.Println(lc.logLevel, lc.fileName, lc.maxSize, lc.maxBackups, lc.maxAge, lc.async, lc.cacheMaxCount, lc.batchSize)
	lc.Level = strings.ToLower(lc.Level)
	if paramsCK := func() error {
		if lc.Level != "info" && lc.Level != "debug" && lc.Level != "warn" && lc.Level != "error" {
			return errors.New("params is illegal")
		} else {
			return nil
		}
	}(); paramsCK != nil {
		return nil, paramsCK
	}
	if lc.CallerSkip == 0{
		lc.CallerSkip = 1
	}
	userPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= func() zapcore.Level {
			switch lc.Level {
			case "debug":
				{
					return zapcore.DebugLevel
				}
			case "info":
				{
					return zapcore.InfoLevel
				}
			case "warn":
				{
					return zapcore.WarnLevel
				}
			case "error":
				{
					return zapcore.ErrorLevel
				}
			default:
				{
					return zapcore.ErrorLevel
				}
			}
		}()
	})
	lumberjackInst := &lumberjack.Logger{
		Filename:   lc.File,
		MaxSize:    lc.MaxSize, // megabytes
		MaxBackups: lc.MaxBackup,
		MaxAge:     lc.MaxAge, // days

		Async:         lc.Async,
		CacheMaxCount: lc.CacheMaxCount,
		BatchSize:     lc.BatchSize,
		Wash:          lc.Wash,
	}
	lumberjackInst.Start()
	logRotateUserWriter := zapcore.AddSync(lumberjackInst)
	commonEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	commonEncoder := zapcore.NewJSONEncoder(commonEncoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(commonEncoder, logRotateUserWriter, userPriority),
	)
	if lc.Caller {
		localLog = &Logger{lumberjack: lumberjackInst, logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(lc.CallerSkip), zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: globalpid})).Sugar()}
	} else {
		localLog = &Logger{lumberjack: lumberjackInst, logger: zap.New(core, zap.AddCallerSkip(1), zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: globalpid})).Sugar()}
	}
	return localLog, nil
}


func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)

}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Warn( keysAndValues ...interface{}) {
	l.logger.Warn(keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *Logger)Debug(args ...interface{}){
	l.logger.Debug(args...)
}

func (l *Logger)Info(args ...interface{}){
	l.logger.Info(args...)
}

func (l *Logger)Error(args ...interface{}){
	l.logger.Error(args...)
}


func (l *Logger) Close() {
	l.lumberjack.Stop()
}

