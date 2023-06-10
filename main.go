package main

import (
	"github.com/mct-joken/jkojs-agent/client"
	"github.com/mct-joken/jkojs-agent/lib"
)

func main() {
	err := lib.LoadConfig()
	if err != nil {
		panic(err)
	}
	lib.InitLogger()
	client.AutoFetcher()
}
