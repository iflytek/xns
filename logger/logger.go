package logger

import "github.com/xfyun/xns/conf"

type LoggerI interface {
	Debug(args ...interface{})
	Info(kvs ...interface{})
	Warn(args ...interface{})
	Error(kvs ...interface{})
	Infow(f string,args ...interface{})
	Errorw(f string,args ...interface{})
}

var accessLogInst LoggerI // 集群api访问日志

var adminLogInst LoggerI // admin api 操作日志
var (
	runtimeLogInst     LoggerI // 程序运行时日志
	debugLogInst       LoggerI
	errorLogInst       LoggerI
	clusterEventLogger LoggerI // 集群事件更新日志
)

var (
	loggerInstances = map[string]LoggerI{}
)

func Init(logs map[string]*conf.LogConf) (err error) {
	for name, logConf := range logs {
		loggerInstances[name], err = newLogger(logConf)
		if err != nil {
			return err
		}
	}

	accessLogInst = getLogger("access")
	adminLogInst = getLogger("admin")
	runtimeLogInst = getLogger("runtime")
	debugLogInst = getLogger("debug")
	errorLogInst = getLogger("error")
	clusterEventLogger = getLogger("cluster")
	return nil
}

func newLogger(conf *conf.LogConf) (LoggerI, error) {
	//todo
	return newLocalLog(conf)
	//return nil, nil
}

func getLogger(name string) LoggerI {
	logger, ok := loggerInstances[name]
	if !ok {
		panic("logger name '" + name + "' not defined in config loggers")
	}
	return logger
}

func Access() LoggerI {
	return accessLogInst
}

func Debug() LoggerI {
	return debugLogInst
}

func Err() LoggerI {
	return errorLogInst
}

func Admin() LoggerI {
	return adminLogInst
}

func Runtime() LoggerI {
	return runtimeLogInst
}

func Event() LoggerI {
	return clusterEventLogger
}
