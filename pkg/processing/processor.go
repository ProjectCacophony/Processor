package processing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/Cacophony/Processor/pkg/kit"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

// TODO: async processing
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

	for d := range deliveries {
		var event events.Event
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			p.logger.Error("failed to unmarshal event", zap.Error(err))
		}

		if event.Type == events.MessageCreateEventType {
			if event.MessageCreate.Author.Bot {
				err = d.Ack(false)
				if err != nil {
					p.logger.Error("failed to ack event", zap.Error(err))
				}

				continue
			}

			p.logger.Info(
				fmt.Sprintf("%s: %s", event.MessageCreate.Author.String(), event.MessageCreate.Content),
			)

			if event.MessageCreate.Content == "ping" {
				createdAt, _ := event.MessageCreate.Timestamp.Parse()

				session, _ := kit.BotSession(event.BotUserID)

				_, _ = session.ChannelMessageSend(
					event.MessageCreate.ChannelID,
					"latency\ndiscord to gateway: "+event.ReceivedAt.Sub(createdAt).String()+"\n"+
						"gateway to processor: "+time.Since(event.ReceivedAt).String(),
				)
			}
		}

		err = d.Ack(false)
		if err != nil {
			p.logger.Error("failed to ack event", zap.Error(err))
		}
	}

	return nil
}
