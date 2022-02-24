package pubsub

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

// GroupConsumer represents a Sarama consumer group consumer

type NewGroupTaskFunc func(*Message) (GroupTask, error)

type GroupConsumer struct {
	Ready chan bool

	ch      chan GroupTask
	newFunc NewGroupTaskFunc

	wg     sync.WaitGroup
	client sarama.ConsumerGroup
	cancel context.CancelFunc
}

type GroupTask interface {
	Finish()
	Wait()
}

type Message struct {
	Lag int64
	*sarama.ConsumerMessage
}

func NewGroupConsumer(client sarama.ConsumerGroup, newFunc NewGroupTaskFunc) *GroupConsumer {
	return &GroupConsumer{
		Ready:   make(chan bool),
		ch:      make(chan GroupTask),
		newFunc: newFunc,
		client:  client,
	}
}

func (consumer *GroupConsumer) Output() <-chan GroupTask {
	return consumer.ch
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *GroupConsumer) Setup(session sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	log.Printf("Setup session : %v", session)
	//session.ResetOffset("t2p4", 0, 13, "")
	//GlobalApp.logger..Info(session.Claims())
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *GroupConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("Cleanup session : %v", session)
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *GroupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		//consumer.logger.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		lag := claim.HighWaterMarkOffset() - (message.Offset)

		m, err := consumer.newFunc(&Message{lag, message})
		if err != nil {
			log.Printf("newFunc(): err %v", err)
			continue
		}
		consumer.ch <- m
		m.Wait()

		session.MarkMessage(message, "")
		// consumer.logger.Debugf("session.MarkMessage ")
		session.Commit()
		// consumer.logger.Debugf("session.Commit ")
	}
	return nil
}

func (consumer *GroupConsumer) StartConsume(gctx context.Context, allTopic string, consumerCount int) error {
	//消费者数量
	client := consumer.client
	topics := strings.Split(allTopic, ",")

	ctx, cancel := context.WithCancel(gctx)
	consumer.cancel = cancel

	if consumerCount == 0 {
		consumerCount = 1
	}
	consumer.wg.Add(consumerCount)
	go func() {
		defer consumer.wg.Done()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Printf("Sarama %v consumer up and running!...", topics)
	return nil
}

func (consumer *GroupConsumer) StopConsume() error {
	defer close(consumer.ch)

	consumer.cancel()
	consumer.wg.Wait()

	if err := consumer.client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
		return err
	}
	return nil
}
