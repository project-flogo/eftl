package examples

import (
	"context"
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
	"github.com/project-flogo/eftl/lib"
	"github.com/project-flogo/microgateway/api"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	go StartFTL()
	time.Sleep(10 * time.Second)
	go StartEFTL()
	time.Sleep(10 * time.Second)
	os.Exit(m.Run())
}

type handler struct {
	Hit bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeHttp")
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
	h.Hit = true
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
	testHandler := handler{}
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
	for !testHandler.Hit {

	}
}

func TestIntegrationAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	e, err := ConsumerExample()
	assert.Nil(t, err)
	testConsumer(t, e)
}

func TestIntegrationJSON(t *testing.T) {
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
