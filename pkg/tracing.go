package pkg

type Tracing struct {
	Http Http `yaml:"http" json:"http,omitempty"`
}

type Http struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
}
