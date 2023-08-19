package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	"net/http"
	"os"
	"os/signal"

	"github.com/jhekau/gdown/pkg/core/models/logger"
)

const defStopCode = http.StatusInternalServerError
const defSignalMax = 1

type HTTP struct{

	l logger.Logger

	// по умолчанию 1 (defSignalMax)
	// максимальное кол-во сигналов о завершении работы
	// после которого, будет быстрое завершение программы c паникой без 
	// ожидания обработки запросов другими потоками. 
	incSignalMax int32

	// по умолчанию 500 (defStopCode), когда получен сигнал о завершении работы, отправляется
	// httpstatus code клиенту, чтобы он принял решение о дальнейшей судьбе запроса. 
	// Например, nginx может быть настроен так, что при получении http code 500, все новые
	// запросы будет перенаправлять на другие воркеры.
	httpStopCode int

	serv *http.Server
	sig chan os.Signal

	sCtrl *signalControl
	cCtrl *connectControl
}

func (h *HTTP) signalInit(signals ...os.Signal) {

	h.sig = make(chan os.Signal, 1)
	signal.Notify(h.sig, signals...)

	h.sCtrl = &signalControl{h: h}
	h.sCtrl.start()
}

func (h *HTTP) newServerWithHandler( fn http.HandlerFunc, ) *http.Server {

	h.serv = &http.Server{}
	
	// connection control
	h.cCtrl = &connectControl{}
	h.serv.ConnState = h.cCtrl.serverOnStateChange

	h.serv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// conncetion control
		if !h.cCtrl.newReq() {
			w.WriteHeader(h.httpStopCode)
			return
		}

		fn(w, r)
	})

	return h.serv
}






