package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	"net/http"
	"syscall"

	logger_default "github.com/jhekau/gdown/internal/core/logs/default"
)

func NewServerWithHandler( fn http.HandlerFunc ) (*http.Server, *settings) {

	h := &HTTP{
		l : &logger_default.Logger{},
	}

	s := &settings{h}
	s.setDefault()

	h.signalInit(
		syscall.SIGKILL,
		syscall.SIGTERM, // pid package
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	return h.newServerWithHandler(fn), s
}