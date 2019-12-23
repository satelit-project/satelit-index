package main

import (
	"os"
	"os/signal"
	"syscall"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/indexing/anidb"
	"shitty.moe/satelit-project/satelit-index/logging"
	"shitty.moe/satelit-project/satelit-index/server"
	"shitty.moe/satelit-project/satelit-index/task"
)

func main() {
	log, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	cfg := readConfig(log)
	q := makeQueries(cfg, log)
	shed := makeTaskScheduler(cfg, q, log)
	srv := server.New(cfg, q, log)

	shed.Start()
	defer shed.Stop()

	if err = srv.Run(); err != nil {
		log.Fatalf("error while serving files: %v", err)
	}

	log.Infof("server started")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Infof("stopping server")

	if err := srv.Shutdown(); err != nil {
		log.Errorf("failed to shutdown server: %v", err)
		return
	}
}

func readConfig(log *logging.Logger) config.Config {
	cfg, err := config.Default()
	if err != nil {
		log.Fatalf("failed to read app configuration: %v", err)
	}

	return cfg
}

func makeQueries(cfg config.Config, log *logging.Logger) *db.Queries {
	dbf := db.NewFactory(cfg.Database, log)
	q, err := dbf.MakeQueries()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	return q
}

func makeTaskScheduler(cfg config.Config, q *db.Queries, log *logging.Logger) task.Scheduler {
	sh := task.NewScheduler(log)

	upd := anidb.IndexUpdateTaskFactory{
		Cfg: cfg.AniDB,
		DB:  q,
		Log: log,
	}
	sh.Add(upd)

	return sh
}
