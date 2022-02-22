package pubsub

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wjiec/gdsn"
	"github.com/xdg-go/scram"
)

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func CreateKafkaSubscriber(kafkaDSN string, group string) (sarama.ConsumerGroup, error) {
	d, err := gdsn.Parse(kafkaDSN)
	if err != nil {
		return nil, err
	}
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionGZIP

	config.Version = sarama.V2_7_0_0
	config.Metadata.Retry.Backoff = time.Second * 5
	config.Consumer.Return.Errors = true

	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.Session.Timeout = time.Second * 6
	config.Consumer.Group.Heartbeat.Interval = time.Second * 2
	config.Consumer.MaxProcessingTime = time.Minute * 10

	if cid := d.Query().Get("ClientID"); cid != "" {
		config.ClientID = cid
	}
	// 账号密码访问
	if d.User.Username() != "" {
		config.Metadata.Full = true
		config.Net.SASL.Enable = true
		config.Net.SASL.User = d.User.Username()
		config.Net.SASL.Password, _ = d.User.Password()

		config.Net.SASL.Handshake = true
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{
				HashGeneratorFcn: SHA256,
			}
		}
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	}

	return sarama.NewConsumerGroup(strings.Split(d.Address(), ","), group, config)
}

func CreateKafkaPublisher(kafkaDSN string) (sarama.SyncProducer, error) {
	d, err := gdsn.Parse(kafkaDSN)
	if err != nil {
		return nil, err
	}
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionGZIP
	config.Version = sarama.V2_7_0_0
	config.Metadata.Retry.Backoff = time.Second * 5
	config.Consumer.Return.Errors = true
	// 账号密码访问
	if d.User.Username() != "" {
		config.Metadata.Full = true
		config.Net.SASL.Enable = true
		config.Net.SASL.User = d.User.Username()
		config.Net.SASL.Password, _ = d.User.Password()

		config.Net.SASL.Handshake = true
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{
				HashGeneratorFcn: SHA256,
			}
		}
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	}
	return sarama.NewSyncProducer(strings.Split(d.Address(), ","), config)
}
