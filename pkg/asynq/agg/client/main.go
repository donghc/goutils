package main

import (
	"context"
	"github.com/hibiken/asynq"
	"log"
	"strings"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func aggregate(group string, tasks []*asynq.Task) *asynq.Task {
	log.Printf("Aggregating %d tasks from group %q", len(tasks), group)
	var b strings.Builder
	for _, t := range tasks {
		b.Write(t.Payload())
		b.WriteString("\n")
	}
	return asynq.NewTask("aggregated-task", []byte(b.String()))
}

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
				"tutorial":    1,
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
			GroupAggregator:  asynq.GroupAggregatorFunc(aggregate),
			GroupGracePeriod: 10 * time.Second,
			GroupMaxDelay:    30 * time.Second,
			GroupMaxSize:     20,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("aggregated-task", handleAggregatedTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func handleAggregatedTask(ctx context.Context, task *asynq.Task) error {
	log.Print("Handler received aggregated task")
	log.Printf("aggregated messags:%v, %s", task.Type(), task.Payload())
	return nil
}
