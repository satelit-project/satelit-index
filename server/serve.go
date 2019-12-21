package server

import (
	"fmt"
	"net/http"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

type IndexServer struct {
	Cfg *config.Config
	Log *logging.Logger
}

func (s IndexServer) Run() error {
	dir := http.Dir(s.Cfg.Serving.Path)
	fs := http.FileServer(dir)

	addr := fmt.Sprintf(":%d", s.Cfg.Serving.Port)
	http.Handle("/index/", http.StripPrefix("/index/", fs))
	return http.ListenAndServe(addr, nil)
}
