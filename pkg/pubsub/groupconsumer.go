package pubsub

import (
	"context"
	"strings"
	"sync"

	zaploger "github.com/donghc/goutils/pkg/log"

	"github.com/Shopify/sarama"
)

// GroupConsumer represents a Sarama consumer group consumer

type NewGroupTaskFunc func(*Message) (GroupTask, error)

type GroupConsumer struct {
	Ready  chan bool
	logger zaploger.Logger

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

func NewGroupConsumer(client sarama.ConsumerGroup, newFunc NewGroupTaskFunc,
	logger zaploger.Logger) *GroupConsumer {
	return &GroupConsumer{
		Ready:   make(chan bool),
		logger:  logger,
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
	consumer.logger.Debugf("Setup session : %v", session)
	//session.ResetOffset("t2p4", 0, 13, "")
	//GlobalApp.logger..Info(session.Claims())
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *GroupConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	consumer.logger.Debugf("Cleanup session : %v", session)
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *GroupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		//consumer.logger.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		lag := claim.HighWaterMarkOffset() - (message.Offset)

		m, err := consumer.newFunc(&Message{lag, message})
		if err != nil {
			consumer.logger.Errorf("newFunc(): err", err)
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

func (consumer *GroupConsumer) StartConsume(gctx context.Context, allTopic string) error {
	//消费者数量
	client := consumer.client
	topics := strings.Split(allTopic, ",")
	logger := consumer.logger

	ctx, cancel := context.WithCancel(gctx)
	consumer.cancel = cancel

	consumer.wg.Add(1)
	go func() {
		defer consumer.wg.Done()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				logger.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	logger.Infof("Sarama %v consumer up and running!...", topics)
	return nil
}

func (consumer *GroupConsumer) StopConsume() error {
	defer close(consumer.ch)

	consumer.cancel()
	consumer.wg.Wait()

	if err := consumer.client.Close(); err != nil {
		consumer.logger.Panicf("Error closing client: %v", err)
		return err
	}
	return nil
}
