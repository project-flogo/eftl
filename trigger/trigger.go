package trigger

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/eftl/lib"
)

var triggerMetadata = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

// Handler is a EFTL handler
type Handler struct {
	handler  trigger.Handler
	settings HandlerSettings
}

// Trigger is a simple EFTL trigger
type Trigger struct {
	connection *lib.Connection
	stop       chan bool
	config     *trigger.Config
	settings   *Settings
	handlers   map[string]Handler
	logger     log.Logger
}

// Factory is a EFTL trigger factory
type Factory struct {
}

// New creates a new EFTL trigger
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{config: config, settings: s}, nil
}

// Metadata returns the EFTL trigger metadata
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMetadata
}

// Initialize initializes the trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	logger := ctx.Logger()
	t.logger = logger

	handlers := make(map[string]Handler)
	for _, handler := range ctx.GetHandlers() {
		settings := HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), &settings, true)
		if err != nil {
			return err
		}

		handlers[settings.Dest] = Handler{
			handler:  handler,
			settings: settings,
		}
	}
	t.handlers = handlers

	return nil
}

// Start implements trigger start
func (t *Trigger) Start() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	ca := t.settings.CA
	if ca != "" {
		certificate, err := ioutil.ReadFile(ca)
		if err != nil {
			t.logger.Error("can't open certificate", err)
			return err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificate)
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
	}

	id := t.settings.ID
	user := t.settings.User
	password := t.settings.Password
	options := &lib.Options{
		ClientID:  id,
		Username:  user,
		Password:  password,
		TLSConfig: tlsConfig,
	}

	url := t.settings.URL
	errorsChannel := make(chan error, 1)
	var err error
	t.connection, err = lib.Connect(url, options, errorsChannel)
	if err != nil {
		t.logger.Errorf("connection failed: %s", err)
		return err
	}

	messages := make(chan lib.Message, 1000)
	for dest := range t.handlers {
		matcher := fmt.Sprintf("{\"_dest\":\"%s\"}", dest)
		_, err = t.connection.Subscribe(matcher, "", messages)
		if err != nil {
			t.logger.Errorf("subscription failed: %s", err)
			return err
		}
	}

	t.stop = make(chan bool, 1)
	go func() {
		for {
			select {
			case message := <-messages:
				value := message["_dest"]
				dest, ok := value.(string)
				if !ok {
					t.logger.Errorf("dest is required for valid message")
					continue
				}
				handler, ok := t.handlers[dest]
				if !ok {
					t.logger.Error("no handler for dest ", dest)
					continue
				}
				value = message["content"]
				content, ok := value.([]byte)
				if !ok {
					content = []byte{}
				}
				err := t.RunAction(handler, dest, content)
				if err != nil {
					t.logger.Errorf("action error: %s", err)
				}
			case err := <-errorsChannel:
				t.logger.Errorf("connection error: %s", err)
			case <-t.stop:
				return
			}
		}
	}()

	return nil
}

// Stop implements trigger stop
func (t *Trigger) Stop() error {
	if t.connection != nil {
		t.connection.Disconnect()
	}
	if t.stop != nil {
		t.stop <- true
	}
	return nil
}

// RunAction starts a new Process Instance
func (t *Trigger) RunAction(handler Handler, dest string, content []byte) error {
	t.logger.Infof("EFTL Trigger: Received request for id '%s'", t.config.Id)

	var data map[string]interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	replyTo := ""
	if value, ok := data["replyTo"]; ok && value != "" {
		replyTo = value.(string)
	}

	output := Output{
		Content: data,
	}

	results, err := handler.handler.Handle(context.Background(), &output)
	if err != nil {
		return err
	}

	if replyTo == "" {
		return nil
	}

	reply := &Reply{}
	reply.FromMap(results)

	return t.connection.Publish(lib.Message{
		"_dest":   replyTo,
		"content": reply.Data,
	})
}
