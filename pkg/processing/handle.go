package processing

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/Cacophony/Processor/plugins"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Processor) handle(delivery amqp.Delivery) error {
	var event events.Event
	err := json.Unmarshal(delivery.Body, &event)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.processingDeadline)
	defer cancel()

	event.WithContext(ctx)
	event.WithTokens(p.discordTokens)
	event.WithLocalisations(plugins.LocalisationsList)
	event.WithState(p.stateClient)
	event.WithBotOwnerIDs(p.botOwnerIDs)
	event.WithPaginator(p.paginator)
	event.WithLogger(p.logger.With(zap.String("service", "processor")))
	event.WithRedis(p.redis)
	event.WithDB(p.db)

	event.Parse()

	err = delivery.Ack(false)
	if err != nil {
		return errors.Wrap(err, "failed to ack event")
	}

	switch event.Type {
	case events.MessageCreateType:
		err = p.paginator.HandleMessageCreate(event.MessageCreate)
		if err != nil {
			event.ExceptSilent(err)
		}
	case events.MessageReactionAddType:
		err = p.paginator.HandleMessageReactionAdd(event.MessageReactionAdd)
		if err != nil {
			event.ExceptSilent(err)
		}
	}

	var wg sync.WaitGroup
	var handled bool
	for _, plugin := range plugins.PluginList {
		if !p.featureFlagger.IsEnabled(featureFlagPluginKey(plugin.Name()), true) {
			p.logger.Debug("skipping plugin as it is disabled by feature flags",
				zap.String("plugin_name", plugin.Name()),
			)
			continue
		}

		event.WithLogger(p.logger.With(zap.String("plugin", plugin.Name())))

		if plugin.Passthrough() {
			// if passthrough, continue with next plugin asap

			wg.Add(1)

			go func(pl plugins.Plugin) {
				defer wg.Done()

				p.executePlugin(pl, &event)
			}(plugin)
			continue
		}

		// else, wait for result, exit if handled
		handled = p.executePlugin(plugin, &event)
		if handled {
			return nil
		}
	}

	wg.Wait()

	return nil
}

func (p *Processor) executePlugin(plugin plugins.Plugin, event *events.Event) bool {
	defer func() {
		err := recover()
		if err != nil {
			if _, ok := err.(error); ok {
				event.ExceptSilent(err.(error))
			}

			p.logger.Error("plugin failed to handle event",
				zap.String("plugin", plugin.Name()),
				zap.String("event_type", string(event.Type)),
				zap.Any("error", err),
			)
		}
	}()

	return plugin.Action(event)
}

func featureFlagPluginKey(pluginName string) string {
	return "plugin-" + pluginName
}
