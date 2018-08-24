package ping

import (
	"time"

	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/dhelpers"
)

func simplePing(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "ping.simplePing")
	defer span.Finish()

	// send message posting time since we received the event
	_, err := event.SendMessage(event.MessageCreate.ChannelID, time.Since(event.ReceivedAt).String())
	if err != nil {
		panic(err)
	}
}

func pingInfo(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "ping.pingInfo")
	defer span.Finish()

	// create message
	message := "pong, Gateway => Processor: " + time.Since(event.ReceivedAt).String() + "\n"

	// post message
	_, err := event.SendMessage(event.MessageCreate.ChannelID, message)
	if err != nil {
		panic(err)
	}
}
