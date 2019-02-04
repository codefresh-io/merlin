package logger

import (
	"fmt"
	"os"

	"github.com/onrik/logrus/filename"
	"github.com/rifflock/lfshook"
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
		// log.SetFormatter(&logrus.JSONFormatter{})
		log.AddHook(filename.NewHook())
		log.SetLevel(logrus.DebugLevel)
	}
	path := fmt.Sprintf("%s/.merlinlog", os.Getenv("PWD"))
	os.Remove(path)
	log.AddHook(lfshook.NewHook(path, &logrus.JSONFormatter{}))
	requestLogger := log.WithFields(opt.Fields)
	return requestLogger

}
