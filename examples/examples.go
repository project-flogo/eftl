package examples

import (
	"github.com/project-flogo/contrib/activity/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/eftl/trigger"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
)

// ConsumerExample returns an EFTL consumer example
func ConsumerExample() (engine.Engine, error) {
	app := api.NewApp()

	gateway := microapi.New("Pets")
	service := gateway.NewService("PetStorePets", &rest.Activity{})
	service.SetDescription("Get pets by ID from the petstore")
	service.AddSetting("uri", "http://localhost:8181/a")
	service.AddSetting("method", "POST")
	service.AddSetting("headers", map[string]string{
		"Accept": "application/json",
	})
	step := gateway.NewStep(service)
	step.AddInput("content", "=$.payload.content")
	response := gateway.NewResponse(false)
	response.SetCode(200)
	response.SetData("=$.PetStorePets.outputs.data")
	settings, err := gateway.AddResource(app)
	if err != nil {
		panic(err)
	}

	trg := app.NewTrigger(&trigger.Trigger{}, &trigger.Settings{URL: "ws://localhost:9191/channel"})
	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		Dest: "sample",
	})
	if err != nil {
		return nil, err
	}

	_, err = handler.NewAction(&microgateway.Action{}, settings)
	if err != nil {
		return nil, err
	}

	return api.NewEngine(app)
}
