package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)


func Connect(conf Config) (redis.UniversalClient, error) {
	var client redis.UniversalClient
	if len(conf.Cluster) >0 {
		client = redis.NewClusterClient(conf.ClusterOption())
	}else {
		client = redis.NewClient(conf.Options())
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connect redis success:", conf.Host)

	return client, nil
}
