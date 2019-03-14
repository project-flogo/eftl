package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/eftl/docker"
	"github.com/project-flogo/eftl/examples"
	"github.com/project-flogo/eftl/lib"
)

var (
	app    = flag.Bool("app", false, "run the flogo app")
	target = flag.Bool("target", false, "run the target server")
	client = flag.Bool("client", false, "send a message")
)

func main() {
	flag.Parse()

	if *app {
		fmt.Println("Starting EFTL...")
		docker.StartEFTL()
		fmt.Println("EFTL started")
		e, err := examples.ConsumerExample()
		if err != nil {
			panic(err)
		}
		engine.RunEngine(e)
	} else if *target {
		handler := func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			mime := r.Header.Get("Content-Type")
			log.Println(r.RequestURI)
			log.Println(mime)
			log.Println(string(body))
			w.Header().Set("Content-Type", mime)
			_, err = w.Write(body)
			if err != nil {
				log.Println(err)
			}
		}
		http.HandleFunc("/a", handler)
		http.HandleFunc("/b", handler)
		http.HandleFunc("/c", handler)

		log.Fatal(http.ListenAndServe(":8181", nil))
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
		connection.Publish(lib.Message{
			"_dest":   "sample",
			"content": []byte(`{"message": "hello world"}`),
		})
	} else {
		flag.PrintDefaults()
	}
}
