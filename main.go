package main

import (
	"fmt"
	"github.com/mct-joken/jkojs-agent/client"
	"github.com/mct-joken/jkojs-agent/cmd"
	"github.com/mct-joken/jkojs-agent/lib"
	"io"
	"os"
)

func main() {
	err := lib.LoadConfig()
	if err != nil {
		panic(err)
	}
	lib.InitLogger()
	if os.Getenv("DEBUG") == "0" {
		f, err := os.Open(os.Args[3])
		if err != nil {
			return
		}
		file, _ := io.ReadAll(f)
		fmt.Println(os.Args[3], file)
		cmd.Start(os.Args[1], string(file), os.Args[2])
		return
	}
	client.AutoFetcher()
}
