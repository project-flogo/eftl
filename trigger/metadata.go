package trigger

import "github.com/project-flogo/core/data/coerce"

// Settings are the settings of the eftl trigger
type Settings struct {
	URL      string `md:"url,required"`
	ID       string `md:"id"`
	User     string `md:"user"`
	Password string `md:"password"`
	CA       string `md:"ca"`
}

// HandlerSettings are the handler settings of the eftl trriger
type HandlerSettings struct {
	Dest string `md:"dest,required"`
}

// Output are the outputs sent to the action
type Output struct {
	Content interface{} `md:"content"`
}

// FromMap converts the output from map
func (o *Output) FromMap(values map[string]interface{}) error {
	o.Content = values["content"]

	return nil
}

// ToMap converts the outputs to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"content": o.Content,
	}
}

// Reply is the struct used to reply to a request
type Reply struct {
	Code int         `md:"code"`
	Data interface{} `md:"data"`
}

// ToMap converts a repy to a map
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code": r.Code,
		"data": r.Data,
	}
}

// FromMap converts a reply from a map
func (r *Reply) FromMap(values map[string]interface{}) error {
	var err error
	r.Code, err = coerce.ToInt(values["code"])
	if err != nil {
		return err
	}
	r.Data, _ = values["data"]

	return nil
}
