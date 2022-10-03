package api

type Config struct {
	Address         string
	TraceHeader     string
	EnableProfiling bool
}

func (c *Config) New(name string) *Controller {
	return &Controller{name: name, config: c}
}
