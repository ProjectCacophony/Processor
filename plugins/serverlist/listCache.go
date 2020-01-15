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

	messages, err = allChannelMessages(session, channelID)
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

func allChannelMessages(discord *discord.Session, channelID string) ([]*discordgo.Message, error) {
	var messages, partialMessages []*discordgo.Message
	var err error

	var prevBeforeID string
	for {
		partialMessages, err = discord.Client.ChannelMessages(
			channelID,
			100,
			prevBeforeID,
			"",
			"",
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, partialMessages...)

		if len(partialMessages) <= 0 {
			break
		}

		prevBeforeID = partialMessages[len(partialMessages)-1].ID
	}

	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	return messages, nil
}
