package rabbitmq

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"room/infrastructure/config"
	"room/infrastructure/rabbitmq/broker"
)

type ProducerOptions struct {
	ExchangeOpt *ExchangeOption
	BrokerURL   string
}

type Producer struct {
	broker broker.Broker
	cfg    *viper.Viper
	logger *logrus.Logger
}

var producer *Producer
var rankProducer *Producer

func GetProducer() *Producer {

	if producer != nil {
		return producer
	}

	opt := new(ProducerOptions)
	opt.BrokerURL = config.GetConfig().GetString("rabbitMq.url")
	exchangeOpt := &ExchangeOption{
		Name: config.GetConfig().GetString("rabbitMq.exchange"),
		Type: config.GetConfig().GetString("rabbitMq.type"),
	}
	opt.ExchangeOpt = exchangeOpt
	producer = NewProducer(opt)

	return producer
}

func GetRankProducer() *Producer {

	if rankProducer != nil {
		return rankProducer
	}

	opt := new(ProducerOptions)
	opt.BrokerURL = config.GetConfig().GetString("rabbitMq.url")
	exchangeOpt := &ExchangeOption{
		Name: "room-exchange",
		Type: config.GetConfig().GetString("rabbitMq.type"),
	}
	opt.ExchangeOpt = exchangeOpt
	rankProducer = NewProducer(opt)

	return rankProducer
}

func NewProducer(opt *ProducerOptions) *Producer {
	broker := broker.NewAmqpBroker(&broker.AmqpBrokerOptions{
		Url:          opt.BrokerURL,
		Exchange:     opt.ExchangeOpt.Name,
		ExchangeType: opt.ExchangeOpt.Type,
	})
	return &Producer{broker: broker}
}

func (p *Producer) Publish(key string, data interface{}) error {
	var body []byte
	switch d := data.(type) {
	case string:
		body = []byte(d)
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = b
	}
	return p.broker.Publish(key, body)
}

func (p *Producer) PublishDelay(key string, data interface{}, delay int64) error {
	var body []byte
	switch d := data.(type) {
	case string:
		body = []byte(d)
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = b
	}
	return p.broker.PublishDelay(key, body, delay)
}
