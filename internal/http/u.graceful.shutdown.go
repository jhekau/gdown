package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 15 August 2023
 */
import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

const logGS = `github.com/jhekau/gdown/internal/http/graceful.shutdown.go`

func gracefulShutdown(h *HTTP) {

fmt.Println(h, atomic.LoadInt32(&h.sCtrl.c))

	if atomic.LoadInt32(&h.sCtrl.c) > h.incSignalMax {
		h.l.Info(``, `terminating...`)
		log.Fatalln(logGS, `terminating...`)
	}
	atomic.AddInt32(&h.sCtrl.c, 1)
	
	h.l.Info(``, `shutting down...`)
	h.cCtrl.stopWait()

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(h.sCtrl.timeout)*time.Second )

	defer cancel()
	
	if err := h.serv.Shutdown(ctx); err != nil {
		h.l.Error(logGS, "shutdown error:", err)
		defer os.Exit(1)
		return
	}
	defer os.Exit(0)
}