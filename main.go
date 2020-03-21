package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/indexing"
	"shitty.moe/satelit-project/satelit-index/indexing/anidb"
	"shitty.moe/satelit-project/satelit-index/logging"
	"shitty.moe/satelit-project/satelit-index/server"
	"shitty.moe/satelit-project/satelit-index/task"
)

func main() {
	cfg := readConfig()

	log, err := makeLogger(cfg.Logging)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	q := makeQueries(cfg, log)
	storage, err := makeIndexStorage(cfg, log)
	if err != nil {
		log.Errorf("failed to create remote storage: %v", err)
		return
	}

	shed := makeTaskScheduler(cfg, q, storage, log)
	srv, err := server.New(cfg, q, log)
	if err != nil {
		log.Errorf("failed to start server: %v", err)
		return
	}

	shed.Start()
	defer shed.Stop()

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-done

		if err := srv.Shutdown(); err != nil {
			log.Errorf("failed to shutdown server: %v", err)
			return
		}
	}()

	log.Infof("starting server")
	if err = srv.Run(); err != nil {
		log.Errorf("error while serving files: %v", err)
		return
	}

	log.Infof("server stopped")
}

func makeLogger(cfg *config.Logging) (*logging.Logger, error) {
	log, err := logging.NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	if err = log.CaptureSTDLog(); err != nil {
		return nil, err
	}

	return log, nil
}

func readConfig() config.Config {
	cfg, err := config.Default()
	if err != nil {
		panic(fmt.Sprintf("failed to read app configuration: %v", err))
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

func makeIndexStorage(cfg config.Config, log *logging.Logger) (indexing.IndexStorage, error) {
	return indexing.NewIndexStorage(cfg.Storage, cfg.AniDB.StorageDir, log)
}

func makeTaskScheduler(cfg config.Config, q *db.Queries, storage indexing.IndexStorage, log *logging.Logger) task.Scheduler {
	sh := task.NewScheduler(log)

	upd := anidb.IndexUpdateTaskFactory{
		Cfg:     cfg.AniDB,
		DB:      anidb.Queries{Q: q},
		Storage: storage,
		Log:     log,
	}
	sh.Add(upd)

	return sh
}
