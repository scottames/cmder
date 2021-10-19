//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/scottames/cmder/pkg/log"
)

// Example use of the Cmder custom logger

// logger wraps the log.logger to implement custom functionality
type logger struct {
	errLog  log.Logger
	infoLog log.Logger
	warnLog log.Logger
}

// newLogger returns a new logger
func newLogger() *logger {
	return &logger{
		errLog:  log.New().Key("ERR").Color(log.LoggerRed),
		infoLog: log.New().Key("INFO").Color(log.LoggerGreen),
		warnLog: log.New().Key("WARN").Color(log.LoggerYellow),
	}
}

// err logs an error
func (l logger) err(msg string, err ...error) error {
	if len(err) > 0 && err[0] != nil {
		l.errLog.Log(fmt.Errorf("%s %w", msg, err[0]))

		return err[0]
	}

	l.errLog.Log(msg)

	return nil
}

// info logs an informational message
func (l logger) info(msg string) {
	l.infoLog.Log(msg)
}

// warn logs a warning message
func (l logger) warn(msg string) {
	l.warnLog.Log(msg)
}

// mage target

// Loggercustom example custom wrapper using the Cmder logger
func Loggercustom() {
	logger := newLogger()
	logger.info("oh okay...")
	logger.warn("uh oh...")
	logger.err("oh no!")
	logger.err("oh no for real!", fmt.Errorf("this is an error"))
}
