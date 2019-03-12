package examples

import (
	"io"
	"os"
	"os/exec"

	"github.com/project-flogo/contrib/activity/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/eftl/trigger"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
)

// StartFTL starts the ftl service
func StartFTL() {
	cmd := exec.Command("/opt/tibco/ftl/current-version/bin/tibftlserver",
		"--config", "/opt/tibco/eftl/6.0/samples/tibftlserver_eftl.yaml",
		"--name", "SRV1")
	stdout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

// StartEFTL starts the eftl service
func StartEFTL() {
	cmd := exec.Command("/opt/tibco/eftl/6.0/ftl/bin/tibftladmin", "--ftlserver", "http://localhost:8585",
		"--updaterealm", "/opt/tibco/eftl/6.0/samples/tibrealm.json")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("/opt/tibco/eftl/6.0/ftl/bin/tibftlserver",
		"--config", "/opt/tibco/eftl/6.0/samples/tibftlserver_eftl.yaml",
		"--name", "EFTL")
	stdout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

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
