package examples

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/eftl/docker"
	"github.com/project-flogo/eftl/lib"
	"github.com/project-flogo/microgateway/api"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	docker.StartEFTL()
	os.Exit(m.Run())
}

type handler struct {
	Hit chan bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	_, err = w.Write(body)
	if err != nil {
		panic(err)
	}
	h.Hit <- true
}

func Drain(port string) {
	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", port), 10*time.Second)
		if conn != nil {
			conn.Close()
		}
		if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
			break
		}
	}
}

func Pour(port string) {
	for {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", port), 10*time.Second)
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func testConsumer(t *testing.T, e engine.Engine) {
	defer api.ClearResources()

	Drain("8181")
	testHandler := handler{
		Hit: make(chan bool, 1),
	}
	s := &http.Server{
		Addr:           ":8181",
		Handler:        &testHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		s.ListenAndServe()
	}()
	Pour("8181")
	defer s.Shutdown(context.Background())

	err := e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()

	Pour("9191")
	errChannel := make(chan error, 1)
	options := &lib.Options{
		ClientID: "test",
	}
	connection, err := lib.Connect("ws://localhost:9191/channel", options, errChannel)
	if err != nil {
		panic(err)
	}
	defer connection.Disconnect()
	err = connection.Publish(lib.Message{
		"_dest":   "sample",
		"content": []byte(`{"message": "hello world"}`),
	})
	assert.Nil(t, err)

	select {
	case <-testHandler.Hit:
	case <-time.After(30 * time.Second):
		t.Fatal("didn't get message in time")
	}
}

func TestIntegrationConsumerAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	e, err := ConsumerExample()
	assert.Nil(t, err)
	testConsumer(t, e)
}

func TestIntegrationConsumerJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./json/consumer/flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testConsumer(t, e)
}

type Response struct {
	Status string `json:"status"`
}

func testProducer(t *testing.T, e engine.Engine) {
	defer api.ClearResources()

	err := e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()

	Pour("9191")
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
	fmt.Println("there")
	transport := &http.Transport{
		MaxIdleConns: 1,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Transport: transport,
	}
	request := func(payload string) Response {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9096/", bytes.NewReader([]byte(payload)))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		rsp := Response{}
		err = json.Unmarshal(body, &rsp)
		assert.Nil(t, err)
		return rsp
	}
	response := request(`{"message": "hello world"}`)
	assert.Equal(t, response.Status, "Success")

	select {
	case <-messages:
	case <-time.After(30 * time.Second):
		t.Fatal("didn't get message in time")
	}
}

func TestIntegrationProducerAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration producer test in short mode")
	}

	e, err := ProducerExample()
	assert.Nil(t, err)
	testProducer(t, e)
}

func TestIntegrationProducerJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping JSON integration producer test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./json/producer/flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testProducer(t, e)
}
