package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jhekau/gdown/pkg/core/models/logger"
)

const (
	logP = `github.com/jhekau/gdown/internal/http/handle.go`
	logH1 = `H1: `
	logH2 = `H2: `
	logH3 = `H3: `
	logH4 = `H4: `
)

const httpStopCode = 500

type HTTP struct{

	l logger.Logger

	// максимальное кол-во сигналов о завершении работы
	// после которого, будет быстрое завершение программы c паникой без 
	// ожидания обработки запросов другими потоками. 
	incSignalMax int

	// по умолчанию 500, когда получен сигнал о завершении работы, отправляется
	// httpstatus code клиенту, чтобы он принял решение о дальнейшей судьбе запроса. 
	// Например, nginx может быть настроен так, что при получении http code 500, все новые
	// запросы будет перенаправлять на другие воркеры.
	httpStopCode int

	serv *http.Server
	ch chan os.Signal

	sCtrl *sigControl
	cCtrl *connectControl
}

func (h *HTTP) init(signals ...os.Signal) {

	h.ch = make(chan os.Signal, 1)
	signal.Notify(h.ch, signals...)

	// signal control


}

func (h *HTTP) handle( fn http.HandlerFunc, serv *http.Server ) http.HandlerFunc {

	h.serv = serv

	h.init(
		syscall.SIGKILL, syscall.SIGTERM, // pid package
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	// connection control
	h.cCtrl = &connectControl{}
	h.serv.ConnState = h.cCtrl.serverOnStateChange
	
	f := func(w http.ResponseWriter, r *http.Request) {

		// conncetion control
		if !h.cCtrl.newReq() {
			w.WriteHeader(h.httpStopCode)
			return
		}

		fn(w, r)
	}

	return f
}



type sigControl struct {
	h *HTTP
	c int // кол-во полученных сигналов
	cancel chan struct{}
}
func (sc *sigControl) control( sig <-chan os.Signal ) {
	select {
		case <- sc.cancel:
			return
		case <- sig:
			sc.gracefulShutdown()
	}
}
func (sc *sigControl) stop() {
	sc.cancel <- struct{}{}
} 
func (sc *sigControl) gracefulShutdown() {

	sc.h.l.Info(``, `shutting down...`)
	if sc.c++; sc.c > sc.h.incSignalMax {
		log.Fatal(sc.h.l.Info(``, `terminating...`))
	}
	sc.h.cCtrl.stopWait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := sc.h.serv.Shutdown(ctx); err != nil {
		sc.h.l.Error(logP, "shutdown error:", err)
		defer os.Exit(1)
		return
	} else {
		sc.h.l.Info(``, `gracefully stopped`)
	}
	defer os.Exit(0)
}





type connectControl struct{
	c int64 	// кол-во обрабатываемых входящих соединений
	stop int32 	// 1 - больше не принимаем входящие соединения
}

// при приёме запроса, разрешаем или запрещаем его обработку
func (c *connectControl) newReq() bool {
	return !(atomic.LoadInt32(&c.stop) == 1)
}

// икрементирует или дикрементирует кол-во запросов
func (c *connectControl) serverOnStateChange(conn net.Conn, state http.ConnState) {
    switch state {
    case http.StateNew:
        atomic.AddInt64(&c.c, +1)
    case http.StateHijacked, http.StateClosed:
        atomic.AddInt64(&c.c, -1)
    }
}

// не принимаем больше новых подключений и ожидаем окончания всех обрабатываемых
func (c *connectControl) stopWait() {
	atomic.StoreInt32(&c.stop, 1)
	for {
		time.Sleep(1 * time.Second)
		if atomic.LoadInt64(&c.c) == 0 {
			return
		}
	}
}





type settings struct{
	h *HTTP
}

// установка кастомизированныъ системных сигналов, после которых
// http сервер перестанет принимать новые запросы
func (s settings) SignalNotify(signals ...os.Signal) {
	signal.Reset(s.h.ch)
	s.h.init(signals...)
}



