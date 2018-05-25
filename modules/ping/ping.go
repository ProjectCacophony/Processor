package ping

import (
	"time"

	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/dhelpers"
)

func simplePing(ctx context.Context, eventReceivedAt time.Time) {
	// start tracing span
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "ping.simplePing")
	defer span.Finish()

	event := dhelpers.EventFromContext(ctx)

	_, err := event.SendMessage(event.MessageCreate.ChannelID, time.Since(eventReceivedAt).String())
	if err != nil {
		panic(err)
	}
}

func pingInfo(ctx context.Context) {
	// start tracing span
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "ping.pingInfo")
	defer span.Finish()

	event := dhelpers.EventFromContext(ctx)

	message := "pong, Gateway => SqsProcessor: " + time.Since(event.ReceivedAt).String() + "\n"

	_, err := event.SendMessage(event.MessageCreate.ChannelID, message)
	if err != nil {
		panic(err)
	}
}
