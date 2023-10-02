package services

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type singleRedisInstance struct {
	Conn *redis.Client
}

var redisClient *singleRedisInstance
var initOnce sync.Once

func GetRedisInstance() *redis.Client {
	initOnce.Do(func() {
		redisClient = &singleRedisInstance{
			Conn: redis.NewClient(&redis.Options{
				Network:  "tcp",
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		}

		//- check connection after create new client
		_, err := redisClient.Conn.Ping(ctx).Result()
		//- if failed
		if err != nil {
			
			//- this error will be effected of the flow of redis connection => fatal error
			ctx.Done()
			panic(err)
		}

	})
	return redisClient.Conn
}
