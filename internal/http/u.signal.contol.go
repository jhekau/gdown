package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */

type signalControl struct {
	h *HTTP
	c int32 // кол-во полученных сигналов о остановке приложения
	cancel chan struct{}
}

func (sc *signalControl) start() {
	sc.cancel = make(chan struct{})
}
func (sc *signalControl) stop() {
	sc.cancel <- struct{}{}
} 

func (sc *signalControl) control() {
	go func(){
		select {
			case <- sc.cancel:
				return
			case <- sc.h.sig:
				gracefulShutdown(sc.h)
		}
	}()
}

