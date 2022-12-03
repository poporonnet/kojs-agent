package main

import (
	"github.com/mct-joken/jkojs-agent/docker"
	"github.com/mct-joken/jkojs-agent/types"
)

func main() {
	docker.Exec(types.StartExecRequest{})
	// api.StartServer()
}
