package processing

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// TODO: max in flight
// TODO: shutdown logic

// Processor processes incoming events
type Processor struct {
	logger           *zap.Logger
	serviceName      string
	amqpConnection   *amqp.Connection
	amqpExchangeName string
	amqpRoutingKey   string

	amqpChannel *amqp.Channel
	amqpQueue   *amqp.Queue
}

// NewProcessor creates a new processor
func NewProcessor(
	logger *zap.Logger,
	serviceName string,
	amqpConnection *amqp.Connection,
	amqpExchangeName string,
	amqpRoutingKey string,
) (*Processor, error) {
	processor := &Processor{
		logger:           logger,
		serviceName:      serviceName,
		amqpConnection:   amqpConnection,
		amqpExchangeName: amqpExchangeName,
		amqpRoutingKey:   amqpRoutingKey,
	}

	err := processor.init()
	if err != nil {
		return nil, err
	}

	return processor, nil
}

// init declares the exchange, the queue, and the queue binding
func (p *Processor) init() error {
	var err error
	p.amqpChannel, err = p.amqpConnection.Channel()
	if err != nil {
		return errors.Wrap(err, "cannot open channel")
	}

	err = p.amqpChannel.ExchangeDeclare(
		p.amqpExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare exchange")
	}

	ampqQueue, err := p.amqpChannel.QueueDeclare(
		p.serviceName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare queue")
	}
	p.amqpQueue = &ampqQueue

	err = p.amqpChannel.QueueBind(
		p.amqpQueue.Name,
		p.amqpRoutingKey,
		p.amqpExchangeName,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}

	return nil
}

// Start starts processing events
func (p *Processor) Start() error {
	deliveries, err := p.amqpChannel.Consume(
		p.amqpQueue.Name,
		p.serviceName+time.Now().Format(time.RFC3339), // TODO: different consumer tag?
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot consume queue")
	}

	for delivery := range deliveries {
		go func(d amqp.Delivery) {
			err := p.handle(d)
			if err != nil {
				p.logger.Error("failed to handle event",
					zap.Error(err),
				)
			}
		}(delivery)
	}

	return nil
}
