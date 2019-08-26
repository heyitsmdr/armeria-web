package main

import (
	"armeria/internal/pkg/armeria"
	"flag"
)

func main() {
	configPath := flag.String("config", "./config/development.yml", "path to the config file")
	migrateFlag := flag.Bool("migrate", false, "perform a schema migration")

	flag.Parse()

	if *migrateFlag {
		armeria.Init(*configPath, false)
		armeria.Migrate()
	} else {
		armeria.Init(*configPath, true)
	}

}
