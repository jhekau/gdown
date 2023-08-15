package http

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 15 August 2023
 */
import (
	"context"
	"log"
	"os"
	"sync/atomic"
)

const logGS = `github.com/jhekau/gdown/internal/http/graceful.shutdown.go`

func gracefulShutdown(h *HTTP) {

	h.l.Info(``, `shutting down...`)

	if atomic.LoadInt32(&h.sCtrl.c) > h.incSignalMax {
		log.Fatal(h.l.Info(``, `terminating...`))
	}
	atomic.AddInt32(&h.sCtrl.c, 1)
	h.cCtrl.stopWait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := h.serv.Shutdown(ctx); err != nil {
		h.l.Error(logGS, "shutdown error:", err)
		defer os.Exit(1)
		return
	}
	h.l.Info(``, `gracefully stopped`)
	defer os.Exit(0)
}