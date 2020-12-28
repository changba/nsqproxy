package httper

import (
	_ "github.com/ChangbaServer/nsqproxy/internal/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
	_ "net/http/pprof"
)

type Httper struct {
	addr     string
	server   *http.Server
	statikFS http.FileSystem
}

func NewHttper(addr string) *Httper {
	statikFS, err := fs.New()
	if err != nil {
		panic("NewHttper statikFS error: " + err.Error())
	}
	return &Httper{
		addr:     addr,
		server:   &http.Server{Addr: addr, Handler: nil},
		statikFS: statikFS,
	}
}

// 启动HTTP
func (h *Httper) Run() {
	h.router()
	go func() {
		err := h.server.ListenAndServe()
		if err != nil {
			panic("ListenAndServe error: " + err.Error())
		}
	}()
}
