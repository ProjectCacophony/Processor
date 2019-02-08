package processing

import (
	"encoding/json"

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

	err = delivery.Ack(false)
	if err != nil {
		return errors.Wrap(err, "failed to ack event")
	}

	var handled bool
	for _, plugin := range plugins.PluginList {
		if plugin.Passthrough() {
			// if passthrough, continue with next plugin asap

			go p.executePlugin(plugin, event)
			continue
		}

		// else, wait for result, exit if handled
		handled = p.executePlugin(plugin, event)
		if handled {
			return nil
		}
	}

	return nil
}

func (p *Processor) executePlugin(plugin plugins.Plugin, event events.Event) bool {
	defer func() {
		err := recover()
		if err != nil {
			p.logger.Error("plugin failed to handle event",
				zap.String("plugin", plugin.Name()),
				zap.String("event_type", string(event.Type)),
			)
		}
	}()

	return plugin.Action(event)
}
