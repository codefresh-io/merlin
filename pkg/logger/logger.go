package logger

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

type (
	LoggerOptions struct {
		Debug  bool
		Fields map[string]interface{}
	}
)

func New(opt *LoggerOptions) *logrus.Entry {
	log := logrus.New()
	if opt.Debug {
		log.AddHook(filename.NewHook())
		log.SetLevel(logrus.DebugLevel)
	}
	requestLogger := log.WithFields(opt.Fields)
	return requestLogger

}
