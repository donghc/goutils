package main

import (
	"encoding/json"
	"github.com/donghc/goutils/pkg/asynq/task"
	"github.com/hibiken/asynq"
	"log"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     redisAddr,
		Username: "",
		Password: "123456",
	})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	for i := 0; i < 200; i++ {
		task, err := task.NewEmailDeliveryTask(i, "some:template:id")
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		info, err := client.Enqueue(task)
		if err != nil {
			log.Fatalf("创建任务失败: %v", err)
		}
		marshal, err := json.Marshal(info)
		log.Printf("创建任务 task: id=%s queue=%s 成功", info.ID, info.Queue)
		log.Printf("%v", string(marshal))
	}

}
