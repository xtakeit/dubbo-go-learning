package component

import (
	"fmt"

	"go-server/library/clean"
	"go-server/library/http"
)

var HttpServer *http.Server

func SetupHttpServer(port int) (err error) {
	HttpServer = http.NewServer(&http.ServerConf{
		Addr: fmt.Sprintf(":%d", port),
	})

	clean.Push(HttpServer)

	return
}
