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

func (c *Config) New() *Controller {
	return &Controller{
		name:   "job",
		config: c,
	}
}
