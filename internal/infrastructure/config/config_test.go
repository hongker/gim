package config

import (
	"fmt"
	"testing"
)

func TestConfigLoadFile(t *testing.T)   {
	c := New()
	err := c.LoadFile("./test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c.Redis.Port)
	fmt.Println(c.viper.GetInt("redis.port"))
}