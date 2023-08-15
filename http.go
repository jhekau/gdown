package gdown

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import(
	http_ "github.com/jhekau/gdown/internal/http"
)

// graceful shutdown http server
// останавливает приём новых http запросов и ждёт завершения уже обрабатываемых запросов
var HTTPNewServerWithHandler = http_.NewServerWithHandler


