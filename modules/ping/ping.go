package ping

import (
	"time"

	"gitlab.com/Cacophony/dhelpers"
)

func simplePing(event dhelpers.EventContainer, eventReceivedAt time.Time) {
	_, err := event.SendMessage(event.MessageCreate.ChannelID, time.Since(eventReceivedAt).String())
	if err != nil {
		panic(err)
	}
}

func pingInfo(event dhelpers.EventContainer) {
	message := "pong, Gateway => SqsProcessor: " + time.Since(event.ReceivedAt).String() + "\n"

	_, err := event.SendMessage(event.MessageCreate.ChannelID, message)
	if err != nil {
		panic(err)
	}
}
