package main

import (
	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
	"shitty.moe/satelit-project/satelit-index/server"
)

func main() {
	log, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	cfg, err := config.Default()
	if err != nil {
		log.Fatalf("failed to read app configuration: %v", err)
	}

	srv := server.New(cfg, log)
	if err = srv.Run(); err != nil {
		log.Fatalf("error while serving files: %v", err)
	}
}
