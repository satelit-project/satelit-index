package main

import (
	"github.com/satelit-project/satelit-index/config"
	"github.com/satelit-project/satelit-index/server"
)

func main() {
	srv := server.NewServer(*config.ServerConfig())
	srv.Serve()
}
