package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import ()

type signalControl struct {
	h *HTTP
	c int32 // кол-во полученных сигналов о остановке приложения
	timeout int // second
	cancel chan struct{}
}

func (sc *signalControl) start() {
	sc.cancel = make(chan struct{})
	sc.timeout = 5
	sc.control()
}

func (sc *signalControl) stop() {
	sc.cancel <- struct{}{}
} 

func (sc *signalControl) control() {
	go func(){
		for {
			select {
				case <- sc.cancel:
					return
				case <- sc.h.sig:
					gracefulShutdown(sc.h)
			}
		}
	}()
}

