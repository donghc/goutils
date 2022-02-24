package pubsub

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"testing"
)

//MessageContent 消息内容
type MessageContent struct {
	GID      string                 `json:"gid"`
	CreateAt int64                  `json:"create_at"`
	AggAt    int64                  `json:"agg_at"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
}

type Task struct {
	done chan struct{}
}

func (t *Task) Finish() {
	close(t.done)
}

func (t *Task) Wait() {
	<-t.done
}

type ToMessageContent interface {
	Value() *MessageContent
}

type MessageContentTask struct {
	*MessageContent
	*Task
}

func (mct *MessageContentTask) Value() *MessageContent {
	return mct.MessageContent
}

func newDirGroupTask(m *Message) (GroupTask, error) {
	var t MessageContent
	if err := json.Unmarshal(m.Value, &t); err != nil {
		return nil, err
	}
	return &MessageContentTask{&t, &Task{make(chan struct{})}}, nil
}

func procTask(msg *MessageContent) error {
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Println(string(marshal))
	return nil
}

func TestGroupConsumer_StartConsume(t *testing.T) {
	var (
		err      error
		subGroup sarama.ConsumerGroup
		sub      *GroupConsumer
	)

	if subGroup, err = CreateKafkaSubscriber("", ""); err != nil {
		t.Fatalf("CreateKafkaSubscriber err: %v", err)
	}
	sub = NewGroupConsumer(subGroup, newDirGroupTask)

	sub.StartConsume(context.Background(), "", 1)
LOOP:
	for {
		m, ok := <-sub.Output()
		if !ok {
			log.Println("Top file break")
			//continue
			break LOOP
		}
		msg := m.(ToMessageContent).Value()
		if err := procTask(msg); err == nil {
			m.Finish()
		} else {
			log.Panicf("procDirTask: %v", err)
		}
	}
	// close
	// sub.StopConsume()
}
