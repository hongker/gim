package job

type Config struct {
	Address         string
	TraceHeader     string
	EnableProfiling bool
	WorkerNumber    int
}

func (c *Config) New() *Controller {
	return &Controller{
		name:   "job",
		config: c,
	}
}
