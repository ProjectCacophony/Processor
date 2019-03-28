package serverlist

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/discord"
)

func (p *Plugin) getListMessages(
	session *discord.Session,
	channelID string,
) ([]*discordgo.Message, error) {
	data, err := p.redis.Get(redisListMessagesKey(channelID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var messages []*discordgo.Message

	if len(data) > 0 {
		err = json.Unmarshal(data, &messages)
		if err != nil {
			return nil, err
		}

		return messages, nil
	}

	messages, err = session.Client.ChannelMessages( // TODO: query more messages, if required
		channelID,
		100,
		"",
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	err = p.setListMessages(channelID, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (p *Plugin) setListMessages(
	channelID string,
	messages []*discordgo.Message,
) error {
	data, err := json.Marshal(&messages)
	if err != nil {
		return err
	}

	err = p.redis.Set(redisListMessagesKey(channelID), data, 24*time.Hour).Err()
	return err
}
