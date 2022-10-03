package internal

import "testing"

func TestAggregator_Run(t *testing.T) {
	srv := NewConfig().BuildInstance()
	srv.Run()
}
