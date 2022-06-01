package main

import (
	"github.com/donghc/goutils/pkg/asynq/sample/task"
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
				"priority:10": 19,
				"priority:9":  17,
				"priority:8":  15,
				"priority:7":  13,
				"priority:6":  11,
				"priority:5":  9,
				"priority:4":  7,
				"priority:3":  5,
				"priority:2":  3,
				"priority:1":  1,
			},
			StrictPriority: false, //如果是true，则严格按照高中低优先级的顺序执行
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	process := task.NewImageProcessor()
	mux.HandleFunc(task.TypeImageResize, process.ProcessTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
