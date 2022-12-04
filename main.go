package main

import (
	"github.com/mct-joken/jkojs-agent/api"
	"github.com/mct-joken/jkojs-agent/lib"
)

func main() {
	lib.InitLogger()
	api.StartServer()
}
