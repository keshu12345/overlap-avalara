package main

import (
	"flag"
	"log"
	"os"

	"github.com/keshu12345/overlap-avalara/config"
	"github.com/keshu12345/overlap-avalara/internal"
	"github.com/keshu12345/overlap-avalara/logger"
	"github.com/keshu12345/overlap-avalara/server/router"
	"go.uber.org/fx"
)

var configDirPath = flag.String("config", "", "path for config dir")

//var migrationDir = flag.String("migrations", "./migrations", "path for migration dir")

func main() {

	flag.Parse()
	log.New(os.Stdout, "", 0)
	app := fx.New(
		config.NewFxModule(*configDirPath, ""),
		router.Module,
		internal.Module,
		// db.Module,
		// migrate.Module(*migrationDir),
		// dao.Module,
		logger.Module,
	)

	app.Run()
}
