package activity

import "github.com/project-flogo/core/data/coerce"

// Settings are the settings for the EFTL activity
type Settings struct {
	URL      string `md:"url,required"`
	ID       string `md:"id"`
	User     string `md:"user"`
	Password string `md:"password"`
	CA       string `md:"ca"`
}

// Input is the input for the EFTL activity
type Input struct {
	Content interface{} `md:"content,required"`
	Dest    string      `md:"dest,required"`
}

// FromMap converts the Input from a map
func (r *Input) FromMap(values map[string]interface{}) error {
	var err error
	r.Content, err = coerce.ToAny(values["content"])
	if err != nil {
		return err
	}
	r.Dest, err = coerce.ToString(values["dest"])
	if err != nil {
		return err
	}
	return nil
}

// ToMap converts the Input to a map
func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Content": r.Content,
		"Dest":    r.Dest,
	}
}

// Output is the Output of the EFTL activity
type Output struct {
}

// FromMap converts the output from a map
func (o *Output) FromMap(values map[string]interface{}) error {
	return nil
}

// ToMap converts the output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}
