package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	"net"
	"net/http"
	"sync/atomic"
	"time"
)


type connectControl struct{
	c int64 	// кол-во обрабатываемых входящих соединений
	stop int32 	// 1 - больше не принимаем входящие соединения
}

// при приёме запроса, разрешаем или запрещаем его обработку
func (c *connectControl) newReq() bool {
	return !(atomic.LoadInt32(&c.stop) == 1)
}

// икрементирует или дикрементирует кол-во запросов, добавляем в сервер
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