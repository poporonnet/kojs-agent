package main

import (
	"github.com/mct-joken/jkojs-agent/api"
	"github.com/mct-joken/jkojs-agent/lib"
)

func main() {
	err := lib.LoadConfig()
	if err != nil {
		panic(err)
	}
	lib.InitLogger()
	api.StartServer()
}
