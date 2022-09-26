package aggregator

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}
func (c *Config) New() *Aggregator {
	return &Aggregator{}
}
