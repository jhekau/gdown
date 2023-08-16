package logger

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import()

type Logger interface {
	Info(path string, arg ...any)
	Alert(path string, arg ...any)
	Error(path string, arg ...any)
}