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

func (p *Plugin) getChannelServersToPost(
	botID string,
) ([]*ChannelServersToPost, error) {
	data, err := p.redis.Get(channelServersToPost(botID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var serversToPost []*ChannelServersToPost

	err = json.Unmarshal(data, &serversToPost)
	if err != nil {
		return nil, err
	}

	return serversToPost, nil
}

func (p *Plugin) setChannelServersToPost(
	botID string,
	serversToPost []*ChannelServersToPost,
) error {
	data, err := json.Marshal(&serversToPost)
	if err != nil {
		return err
	}

	err = p.redis.Set(channelServersToPost(botID), data, 24*time.Hour*7).Err()
	return err
}
