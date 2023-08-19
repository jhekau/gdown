package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/jhekau/gdown/pkg/core/models/logger"
)

type settings struct{
	h *HTTP
}

// установка кастомизированныъ системных сигналов, после которых
// http сервер перестанет принимать новые запросы
func (s *settings) SignalNotify(signals ...os.Signal) {
	s.h.sCtrl.stop()
	signal.Stop(s.h.sig)
	s.h.signalInit(signals...)
}

// устанавливает максимальное кол-во сигналов о завершении работы
// после которого, будет быстрое завершение программы c паникой без 
// ожидания обработки запросов другими потоками. 
func (s *settings) SignalIncomigMax(c int32) {
	atomic.StoreInt32(&s.h.incSignalMax, c)
}

// устанавливает код ответа клиенту, когда получен сигнал о завершении работы, 
// чтобы клеинт принял решение о дальнейшей судьбе запроса. 
// Например, nginx может быть настроен так, что при получении http code 500, все новые
// запросы будет перенаправлять на другие воркеры. 
func (s *settings) HTTPStopCode(c int) {
	s.h.httpStopCode = c
}

func (s *settings) Logger( l logger.Logger ) {
	s.h.l = l
}

func (s *settings) setDefault() {
	s.h.incSignalMax = defSignalMax
	s.h.httpStopCode = defStopCode
}

func (s *settings) setTimeout(sec int) {
	s.h.sCtrl.timeout = sec
}