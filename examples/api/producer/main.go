package main

import (
	"flag"
	"fmt"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/eftl/examples"
	"github.com/project-flogo/eftl/lib"
)

var (
	ftl    = flag.Bool("ftl", false, "start the ftl server")
	eftl   = flag.Bool("eftl", false, "start the eftl server")
	client = flag.Bool("client", false, "send a message")
	app    = flag.Bool("app", false, "run the flogo app")
)

func main() {
	flag.Parse()

	if *ftl {
		examples.StartFTL()
	} else if *eftl {
		examples.StartEFTL()
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
	} else if *app {
		e, err := examples.ProducerExample()
		if err != nil {
			panic(err)
		}
		engine.RunEngine(e)
	} else {
		flag.PrintDefaults()
	}
}
