package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewGoroutinePool(t *testing.T) {
	gp := NewGoroutinePool(50)
	for {
		gp.Schedule(func() {
			fmt.Println("123")
		})
		time.Sleep(time.Second)
	}

}
