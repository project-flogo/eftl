package main

import (
	"flag"
	"fmt"

	"github.com/project-flogo/eftl/docker"
	"github.com/project-flogo/eftl/lib"
)

var (
	app    = flag.Bool("app", false, "run the flogo app")
	client = flag.Bool("client", false, "send a message")
)

func main() {
	flag.Parse()

	if *app {
		fmt.Println("Starting EFTL...")
		docker.StartEFTL()
		fmt.Println("EFTL started")
	} else if *client {
		errChannel := make(chan error, 1)
		options := &lib.Options{
			ClientID: "test",
		}
		connection, err := lib.Connect("ws://localhost:9191/channel", options, errChannel)
		if err != nil {
			panic(err)
		}
		defer connection.Disconnect()
		messages := make(chan lib.Message, 1000)
		connection.Subscribe(`{"_dest":"sample"}`, "", messages)
		for message := range messages {
			fmt.Println(string(message["content"].([]byte)))
		}
	} else {
		flag.PrintDefaults()
	}
}
