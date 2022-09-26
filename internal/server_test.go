package internal

import "testing"

func TestAggregator_Run(t *testing.T) {
	aggregator := NewConfig().New()
	aggregator.Run()
}
