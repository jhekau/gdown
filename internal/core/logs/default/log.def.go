package logger_default

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import (
	"fmt"
	"log"
)

type Logger struct{}

func (l *Logger) Info(path string, msg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `info`, path, fmt.Sprint(msg...)))
}
func (l *Logger) Warn(path string, msg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `alert`, path, fmt.Sprint(msg...)))
}
func (l *Logger) Error(path string, msg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `error`, path, fmt.Sprint(msg...)))
}

