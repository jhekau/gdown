# Graceful Shutdown
# Demo <!>


### Как использовать?
```Go
import (
    "net/http"
    "github.com/jhekau/gdown"
)

func main(){
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte( `Hello World` ))
    }

    server, _ := gdown.HTTPNewServerWithHandler(handler)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}

```

### Использование альтернативной реализации логгера:
```Go
import (
    "github.com/jhekau/gdown"
    "github.com/jhekau/gdown/pkg/core/models/logger"
)

# check implementation
var _ logger.Logger = (YourLogger)(nil)

server, settings := gdown.HTTPNewServerWithHandler(handler)
settings.Logger(YourLogger)

```

#### Chapters

- v0.0.1: add HTTPNewServerWithHandler;
- v0.0.0: create;
