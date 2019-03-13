package activity

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/eftl/lib"
)

func init() {
	activity.Register(&Activity{}, New)
}

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// New create a New EFTL activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := Settings{}
	err := metadata.MapToStruct(ctx.Settings(), &settings, true)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if settings.CA != "" {
		certificate, err := ioutil.ReadFile(settings.CA)
		if err != nil {
			ctx.Logger().Error("can't open certificate", err)
			return nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificate)
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
	}

	act := Activity{
		options: &lib.Options{
			TLSConfig: tlsConfig,
			ClientID:  settings.ID,
			Username:  settings.User,
			Password:  settings.Password,
		},
		url: settings.URL,
	}

	return &act, nil
}

// Activity is a EFTL client
type Activity struct {
	options *lib.Options
	url     string
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := Input{}
	err = ctx.GetInputObject(&input)
	if err != nil {
		return false, err
	}

	errorsChannel := make(chan error, 1)
	connection, err := lib.Connect(a.url, a.options, errorsChannel)
	if err != nil {
		ctx.Logger().Errorf("connection failed: %s", err)
		return false, err
	}
	defer connection.Disconnect()

	data, err := json.Marshal(input.Content)
	if err != nil {
		ctx.Logger().Errorf("failed to marshal: %s", err)
		return false, err
	}

	err = connection.Publish(lib.Message{
		"_dest":   input.Dest,
		"content": data,
	})
	if err != nil {
		ctx.Logger().Errorf("failed to publish: %s", err)
		return false, err
	}

	return true, nil
}
