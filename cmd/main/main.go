package main

import (
	"flag"

	"github.com/kuzin57/shad-networks/cmd/app"
	"go.uber.org/fx"
)

func main() {
	confPath := flag.String("config", "./config/config.yaml", "path to config")

	fx.New(app.Create(*confPath)).Run()
}
