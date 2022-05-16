package main

import (
	"github.com/donghc/goutils/pkg/asynq/task"
	"github.com/hibiken/asynq"
	"log"
)

const redisAddr = "127.0.0.1:6379"

func main() {

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Network: "tcp", Addr: redisAddr,
			Username: "", Password: "123456",
			DB: 0,
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 0,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: false, //如果是true，则严格按照高中低优先级的顺序执行
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeEmailDelivery, task.HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
