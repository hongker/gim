package api

type Config struct {
	Address         string
	TraceHeader     string
	EnableProfiling bool
}

func (c *Config) New() *Controller {
	return &Controller{name: "api", config: c}
}
