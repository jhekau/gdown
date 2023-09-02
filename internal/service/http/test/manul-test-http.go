package main

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 18 August 2023
 */
import (
	"net/http"

	"github.com/jhekau/gdown"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte( `Hello World` ))
    }

    server, _ := gdown.HTTPNewServerWithHandler(handler)
	server.Addr = `:8080`
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}
