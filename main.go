package main

import (
	"satelit-project/satelit-index/config"
	"satelit-project/satelit-index/cron"
	dbcfg "satelit-project/satelit-index/db"
	"satelit-project/satelit-index/logging"
	"satelit-project/satelit-index/server"

	"github.com/gobuffalo/pop"
)

func main() {
	env := string(config.CurrentEnvironment())
	db, err := pop.Connect(env)
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
	srv.SetLogger(logging.DefaultLogger())
	srv.Serve()
}
