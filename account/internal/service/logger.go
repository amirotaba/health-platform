package service

import (
	"io"
	"log"
	"os"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type logger struct{}

// NewLogger return a Logger.
func NewLogger() domain.Logger {
	return &logger{}
}

// LogError is print messages to log.
func (l *logger) LogError(format string, v ...interface{}) {
	file, err := os.OpenFile("./log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("%s", err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("%s", err)
		}
	}(file)

	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf(format, v...)
}

// LogAccess is print messages to log.
func (l *logger) LogAccess(format string, v ...interface{}) {
	file, err := os.OpenFile("./log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("%s", err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
		}
	}(file)

	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf(format, v...)
}
