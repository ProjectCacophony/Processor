package ping

import (
	"time"

	"context"

	"gitlab.com/Cacophony/dhelpers"
)

func simplePing(ctx context.Context, eventReceivedAt time.Time) {
	event := dhelpers.EventFromContext(ctx)

	_, err := event.SendMessage(event.MessageCreate.ChannelID, time.Since(eventReceivedAt).String())
	if err != nil {
		panic(err)
	}
}

func pingInfo(ctx context.Context) {
	event := dhelpers.EventFromContext(ctx)

	message := "pong, Gateway => SqsProcessor: " + time.Since(event.ReceivedAt).String() + "\n"

	_, err := event.SendMessage(event.MessageCreate.ChannelID, message)
	if err != nil {
		panic(err)
	}
}
