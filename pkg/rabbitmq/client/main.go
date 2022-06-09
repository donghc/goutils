package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// 建立连接
	conn, err := amqp.Dial("amqp://testuser:psKdqSe@192.168.63.34:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 获取channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	var args amqp.Table
	args = amqp.Table{"x-max-priority": int32(10)} // 设置优先级列表的最初优先级
	// 声明队列
	//请注意，我们也在这里声明队列。因为我们可能在发布者之前启动使用者，所以我们希望在尝试使用队列中的消息之前确保队列存在。
	q, err := ch.QueueDeclare(
		"scand_win", // name
		true,        // durable 声明为持久队列
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		args,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.Qos(
		4,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	// 获取接收消息的Delivery通道
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  //  注意这里传false,关闭自动消息确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	batch := make(chan amqp.Delivery, 4)
	ms := make([]amqp.Delivery, 0)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %v: %v: %v", d.Priority, d.MessageId, string(d.Body))
			//batch<-d
			d.Ack(false) // 手动传递消息确认
		}
	}()

	go func() {
	loop:
		for {
			select {
			case m, ok := <-batch:
				if !ok {
					break loop
				}
				log.Println("get message ", fmt.Sprintf("%v", m.MessageId))
				ms = append(ms, m)
				if len(ms) == 4 {
					log.Println("process deal ")
					//time.Sleep(10 * time.Second)
					for _, s := range ms {
						log.Println("ack begin", fmt.Sprintf("%v", s.MessageId))
						s.Ack(false)
						log.Println("ack end", fmt.Sprintf("%v", s.MessageId))
					}
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
