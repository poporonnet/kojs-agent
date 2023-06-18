package main

import (
	"fmt"
	"github.com/mct-joken/jkojs-agent/cmd"
	"github.com/mct-joken/jkojs-agent/pkg/client"
	lib2 "github.com/mct-joken/jkojs-agent/pkg/lib"
	"io"
	"os"
)

func main() {
	err := lib2.LoadConfig()
	if err != nil {
		panic(err)
	}
	lib2.InitLogger()
	if os.Getenv("DEBUG") == "0" {
		f, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		file, _ := io.ReadAll(f)
		fmt.Println(os.Args[2], file)
		cmd.Start(os.Args[3], string(file), os.Args[1])
		return
	}
	client.AutoFetcher()
}
