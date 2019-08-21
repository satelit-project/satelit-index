package main

import (
	"github.com/gobuffalo/pop"
	"github.com/satelit-project/satelit-index/config"
	"github.com/satelit-project/satelit-index/cron"
	dbcfg "github.com/satelit-project/satelit-index/db"
	"github.com/satelit-project/satelit-index/server"
)

func main() {
	db, err := pop.Connect(config.Environment())
	if err != nil {
		panic(err)
	}

	serverCfg := config.ServerConfig()
	anidbCfg := config.AnidbConfig()

	err = dbcfg.SetupAnidbTables(db, anidbCfg)
	if err != nil {
		panic(err)
	}

	anidbJobs := cron.DefaultAnidbScheduler(db, serverCfg, anidbCfg)
	anidbJobs.StartJobs()

	srv := server.NewServer(db, serverCfg)
	srv.Serve()
}
