package logger

import "log"

type stdLogger struct {

}

func (s *stdLogger) Debug(msg string, args ...interface{}) {
	log.Println(msg,args)
}

func (s *stdLogger) Info(msg string, args ...interface{}) {
	log.Println(msg,args)
}

func (s *stdLogger) Warn(msg string, args ...interface{}) {
	log.Println(msg,args)
}

func (s *stdLogger) Error(msg string, args ...interface{}) {
	log.Println(msg,args)
}

