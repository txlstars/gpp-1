package gpplog

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gppio"
	"time"
	"sync"
)

var loggers = make(map[string]*logrus.Logger)

var mutex = sync.Mutex{}

func GetLogger(logName string) *logrus.Logger {
	if logger, ok := loggers[logName]; ok {
		return logger
	}

	mutex.Lock()
	if logger, ok := loggers[logName]; ok{
		return logger
	}
	logger, _ := newLogger("/home/seantian/log", logName)
	loggers[logName] = logger
	mutex.Unlock()

	return logger
}

func newLogger(logPath, logName string) (*logrus.Logger, error) {
	writer, err := rotatelogs.New(
		logPath + "/" + logName + "_%Y%m%d.log",
		// rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(time.Hour * 24),
		rotatelogs.WithMaxAge(time.Hour * 7 * 24),
	)

	if err != nil {
		return nil, err 
	}

	newLog := logrus.New()
	newLog.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel: writer,
			logrus.WarnLevel: writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.JSONFormatter{},
	))

	var nilWriter gppio.EmptyWriter
	newLog.SetOutput(nilWriter)

	return newLog, nil
}
