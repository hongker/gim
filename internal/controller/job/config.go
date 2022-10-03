package job

import "time"

type Config struct {
	Address         string
	TraceHeader     string
	EnableProfiling bool
	WorkerNumber    int

	QueuePollInterval time.Duration
	QueuePollCount    int
}

func (c *Config) New(name string) *Controller {
	return &Controller{
		name:   name,
		config: c,
	}
}
