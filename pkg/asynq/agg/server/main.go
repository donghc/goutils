package main

import (
	"flag"
	"github.com/hibiken/asynq"
	"log"
)

const redisAddr = "127.0.0.1:6379"

var (
	flagRedisAddr = flag.String("redis-addr", "localhost:6379", "Redis server address")
	flagMessage   = flag.String("message", "hello3", "Message to print when task gets processed")
)

func main() {
	flag.Parse()

	c := asynq.NewClient(asynq.RedisClientOpt{Addr: *flagRedisAddr, Password: "123456"})
	defer c.Close()

	task := asynq.NewTask("aggregation-tutorial", []byte(*flagMessage))
	info, err := c.Enqueue(task, asynq.Queue("tutorial"), asynq.Group("example-group"))
	info2, err := c.Enqueue(task, asynq.Queue("priority:10"), asynq.Group("example-group22"))
	if err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}
	log.Printf("Successfully enqueued task: %s", info.ID)
	log.Printf("Successfully enqueued task: %s", info2.ID)
}
