package main

import (
	"go-universal-clipboard/internal/app"
	"go-universal-clipboard/internal/cfg"
)

func main() {
	cfg.LoadEnv(".")
	a := app.NewApp(&cfg.Cfg)
	a.Run()
}
