package processing

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"gitlab.com/Cacophony/go-kit/paginator"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/Cacophony/go-kit/featureflag"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

// TODO: shutdown logic

// Processor processes incoming events
type Processor struct {
	logger                    *zap.Logger
	serviceName               string
	amqpDSN                   string
	amqpExchangeName          string
	amqpRoutingKey            string
	db                        *gorm.DB
	stateClient               *state.State
	concurrentProcessingLimit int
	processingDeadline        time.Duration
	discordTokens             map[string]string
	featureFlagger            *featureflag.FeatureFlagger
	botOwnerIDs               []string
	redis                     *redis.Client

	amqpConnection   *amqp.Connection
	amqpChannel      *amqp.Channel
	amqpQueue        *amqp.Queue
	amqpErrorChannel chan *amqp.Error

	paginator *paginator.Paginator
}

// NewProcessor creates a new processor
func NewProcessor(
	logger *zap.Logger,
	serviceName string,
	amqpDSN string,
	amqpExchangeName string,
	amqpRoutingKey string,
	db *gorm.DB,
	stateClient *state.State,
	concurrentProcessingLimit int,
	processingDeadline time.Duration,
	discordTokens map[string]string,
	featureFlagger *featureflag.FeatureFlagger,
	botOwnerIDs []string,
	redis *redis.Client,
) (*Processor, error) {
	paginator, err := paginator.NewPaginator(
		logger.With(zap.String("feature", "paginator")),
		redis,
		stateClient,
		discordTokens,
	)
	if err != nil {
		return nil, err
	}

	processor := &Processor{
		logger:                    logger,
		serviceName:               serviceName,
		amqpDSN:                   amqpDSN,
		amqpExchangeName:          amqpExchangeName,
		amqpRoutingKey:            amqpRoutingKey,
		db:                        db,
		stateClient:               stateClient,
		concurrentProcessingLimit: concurrentProcessingLimit,
		processingDeadline:        processingDeadline,
		discordTokens:             discordTokens,
		featureFlagger:            featureFlagger,
		botOwnerIDs:               botOwnerIDs,
		redis:                     redis,

		amqpErrorChannel: make(chan *amqp.Error),
		paginator:        paginator,
	}

	err = processor.init()
	if err != nil {
		return nil, err
	}

	return processor, nil
}

func (p *Processor) startErrorHandler() {
	p.amqpConnection.NotifyClose(p.amqpErrorChannel)

	for amqpErr := range p.amqpErrorChannel {

		if amqpErr == nil {
			continue
		}

		if amqpErr.Recover {
			p.logger.Error("received recoverable error from AMQP Broker",
				zap.Any("error", amqpErr),
			)
			continue
		}

		p.logger.Fatal(
			"looks like we lost the connection to the AMQP Broker, will shut down",
			zap.Any("error", amqpErr),
		)
	}

}

// init declares the exchange, the queue, and the queue binding
func (p *Processor) init() error {
	var err error

	p.amqpConnection, err = amqp.Dial(p.amqpDSN)
	if err != nil {
		return errors.Wrap(err, "unable to initialise AMQP session")
	}

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
	go p.startErrorHandler()

	return p.start()
}

func (p *Processor) start() error {
	deliveries, err := p.amqpChannel.Consume(
		p.amqpQueue.Name,
		fmt.Sprintf(
			"%s: launched %s", p.serviceName, time.Now().UTC().Format(time.RFC3339),
		),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot consume queue")
	}

	// keep semaphore channel to limit the amount of events being processed concurrently
	semaphore := make(chan interface{}, p.concurrentProcessingLimit)

	for delivery := range deliveries {
		// wait for channel if channel buffer is full
		semaphore <- nil

		go func(d amqp.Delivery) {
			defer func() {
				// clear channel when completed
				<-semaphore
			}()

			err := p.handle(d)
			if err != nil {
				p.logger.Error("failed to handle event",
					zap.Error(err),
				)
			}
		}(delivery)
	}

	// try to fill channel with amount of buffer size, this makes sure we will wait for all events to finish processing
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- nil
	}

	p.logger.Info("finished Start()")

	return nil
}
