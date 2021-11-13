package serverlist

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleListRefresh(event *events.Event) {
	err := p.clearListCache(event.BotUserID)
	if err != nil {
		event.Except(err)
		return
	}

	err = p.refreshList(event.BotUserID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.list-refresh.content")
	event.Except(err)
}

func (p *Plugin) clearListCache(botID string) error {
	allList, err := p.getList(botID)
	if err != nil {
		return err
	}

	session, err := discord.NewSession(p.tokens, botID)
	if err != nil {
		return fmt.Errorf("failure creating Discord Session for Bot: %w", err)
	}

	for _, item := range allList {
		for _, category := range item.Categories {
			for _, name := range item.Names {

				categoryChannel, err := p.state.Channel(category.Category.ChannelID)
				if err != nil {
					return err
				}

				switch categoryChannel.Type {

				case discordgo.ChannelTypeGuildText:
					// group by not used

					err = p.redis.Del(redisListMessagesKey(categoryChannel.ID)).Err()
					if err != nil {
						return err
					}

				case discordgo.ChannelTypeGuildCategory:

					channel, err := p.getDiscordCategoryChannel(
						session,
						&category.Category,
						name,
					)
					if err != nil {
						return err
					}

					err = p.redis.Del(redisListMessagesKey(channel.ID)).Err()
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (p *Plugin) refreshList(botID string) error {
	allList, err := p.getList(botID)
	if err != nil {
		return fmt.Errorf("failure getting current list: %w", err)
	}

	session, err := discord.NewSession(p.tokens, botID)
	if err != nil {
		return fmt.Errorf("failure creating Discord Session for Bot: %w", err)
	}

	var serversToPost []*ChannelServersToPost

	for _, item := range allList {
		for _, category := range item.Categories {
			for _, name := range item.Names {

				categoryChannel, err := p.state.Channel(category.Category.ChannelID)
				if err != nil {
					return fmt.Errorf("failure getting category channel from state: %w", err)
				}

				switch categoryChannel.Type {

				case discordgo.ChannelTypeGuildText:
					// group by not used

					serversToPost = p.addToServersToPost(
						serversToPost,
						&category.Category,
						categoryChannel,
						*item,
						name,
					)

				case discordgo.ChannelTypeGuildCategory:

					channel, err := p.getDiscordCategoryChannel(
						session,
						&category.Category,
						name,
					)
					if err != nil {
						return fmt.Errorf("failure resolving channel for category channel: %w", err)
					}

					serversToPost = p.addToServersToPost(
						serversToPost,
						&category.Category,
						channel,
						*item,
						name,
					)
				}
			}
		}
	}

	for _, channelToPost := range serversToPost {
		err = p.postChannel(session, channelToPost)
		if err != nil {
			return fmt.Errorf("failure posting channel: %w", err)
		}
	}

	// clear channels if necessary, based on diff with last list
	previousServersToPost, err := p.getChannelServersToPost(botID)
	if err == nil && len(previousServersToPost) > 0 {
		for _, previousChannelToPost := range previousServersToPost {
			if channelServersToPostContain(previousChannelToPost, serversToPost) {
				continue
			}

			err = p.clearListChannel(session, previousChannelToPost)
			if err != nil {
				if errD, ok := err.(*discordgo.RESTError); ok &&
					errD != nil && errD.Message != nil &&
					(errD.Message.Code == discordgo.ErrCodeUnknownChannel ||
						errD.Message.Message == "404: Not Found") {
					continue
				}
				return fmt.Errorf("failure clearing channel list: %w", err)
			}
		}
	}

	err = p.setChannelServersToPost(botID, serversToPost)
	if err != nil {
		return fmt.Errorf("failure updating servers cache: %w", err)
	}

	return nil
}

type ChannelServersToPost struct {
	ChannelID string
	Servers   []*ServerToPost
	Category  *Category
	SortBy    []SortBy
}

type ServerToPost struct {
	Server   *Server
	Name     string
	SortName string
}

func (p *Plugin) clearListChannel(session *discord.Session, channelToPost *ChannelServersToPost) error {
	messages, err := p.getListMessages(session, channelToPost.ChannelID)
	if err != nil {
		return err
	}

	err = discord.DeleteSmart(
		p.redis,
		session,
		channelToPost.ChannelID,
		messages,
		false,
	)
	if err != nil {
		return err
	}

	return p.setListMessages(channelToPost.ChannelID, []*discordgo.Message{})
}

func (p *Plugin) postChannel(session *discord.Session, channelToPost *ChannelServersToPost) error {
	sort.Sort(ServersSorter{
		SortBy:  channelToPost.SortBy,
		Servers: channelToPost.Servers,
	})

	messages, err := p.getListMessages(session, channelToPost.ChannelID)
	if err != nil {
		return err
	}

	var message *discordgo.Message
	var server *ServerToPost
	var i int
	for i, server = range channelToPost.Servers {
		message = nil

		if len(messages) >= i+1 {
			message = messages[i]
		}

		content := p.getMessageContentForServer(server.Server, server.Name)

		if message == nil {
			messages, err = p.newListMessage(
				messages,
				session,
				channelToPost.ChannelID,
				content,
			)
			if err != nil {
				return fmt.Errorf("unable to post list message Channel #%s: %w", channelToPost.ChannelID, err)
			}
		} else if message.Content != content {
			messages, err = p.updateListMessage(
				messages,
				session,
				channelToPost.ChannelID,
				content,
				message.ID,
			)
			if err != nil {
				return fmt.Errorf("unable to update list message Channel #%s Message #%s: %w", channelToPost.ChannelID, message.ID, err)
			}
		}
	}

	if len(messages) >= i+1 {
		for _, message := range messages[i+1:] {
			messages, err = p.deleteListMessage(
				messages,
				session,
				channelToPost.ChannelID,
				message.ID,
			)
			if err != nil {
				var discordError *discordgo.RESTError
				if errors.As(err, &discordError) && discordError.Message != nil && discordError.Message.Code == discordgo.ErrCodeUnknownMessage {
					continue
				}
				return fmt.Errorf("unable to delete list message Channel #%s Message #%s: %w", channelToPost.ChannelID, message.ID, err)
			}
		}
	}

	err = p.setListMessages(channelToPost.ChannelID, messages)
	if err != nil {
		return err
	}

	// update channel topic if needed
	topic := getCategoyTopic(channelToPost.Category, channelToPost.SortBy)
	channel, err := p.state.Channel(channelToPost.ChannelID)
	if err == nil &&
		discord.UserHasPermission(
			p.state, session.BotID, channelToPost.ChannelID, discordgo.PermissionManageChannels,
		) &&
		channel.Topic != topic {
		_, err = session.Client.ChannelEditComplex(channelToPost.ChannelID, &discordgo.ChannelEdit{
			Topic:    topic,
			Position: channel.Position,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) newListMessage(
	messages []*discordgo.Message,
	session *discord.Session,
	channelID string,
	content string,
) ([]*discordgo.Message, error) {
	message, err := discord.SendComplexWithVars(
		session,
		nil,
		channelID,
		&discordgo.MessageSend{
			Content: content,
		},
	)
	if err != nil {
		return nil, err
	}

	messages = append(messages, message[0])

	return messages, nil
}

func (p *Plugin) updateListMessage(
	messages []*discordgo.Message,
	session *discord.Session,
	channelID string,
	content string,
	messageID string,
) ([]*discordgo.Message, error) {
	message, err := discord.EditComplexWithVars(
		p.redis,
		session,
		nil,
		&discordgo.MessageEdit{
			Content: &content,
			ID:      messageID,
			Channel: channelID,
		},
		false,
	)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		if messages[i].ID != message.ID {
			continue
		}

		messages[i] = message
		break
	}

	return messages, nil
}

func (p *Plugin) deleteListMessage(
	messages []*discordgo.Message,
	session *discord.Session,
	channelID string,
	messageID string,
) ([]*discordgo.Message, error) {
	err := discord.Delete(
		p.redis,
		session,
		channelID,
		messageID,
		false,
	)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		if messages[i].ID != messageID {
			continue
		}

		messages = append(messages[:i], messages[i+1:]...)
		break
	}

	return messages, nil
}

func (p *Plugin) getMessageContentForServer(server *Server, name string) string {
	text := fmt.Sprintf(
		"**%s** — https://discord.gg/%s",
		name,
		server.InviteCode,
	)

	if server.Description != "" {
		text += "\n" + server.Description
	}

	return text
}

func (p *Plugin) addToServersToPost(
	serversToPost []*ChannelServersToPost,
	category *Category,
	channel *discordgo.Channel,
	server Server,
	name string,
) []*ChannelServersToPost {
	for i := range serversToPost {
		if serversToPost[i].ChannelID != channel.ID {
			continue
		}

		for j := range serversToPost[i].Servers {
			if serversToPost[i].Servers[j].Server.ID != server.ID {
				continue
			}

			serversToPost[i].Servers[j].Name += "/" + name

			return serversToPost
		}

		serversToPost[i].Servers = append(serversToPost[i].Servers, &ServerToPost{
			Server:   &server,
			Name:     name,
			SortName: getSortName(name),
		})
		return serversToPost
	}

	var sortBy []SortBy

	for _, sortByName := range category.SortBy {
		for _, allSortBy := range allSortBys {
			if string(allSortBy) != sortByName {
				continue
			}

			sortBy = append(sortBy, allSortBy)
		}
	}

	return append(serversToPost, &ChannelServersToPost{
		ChannelID: channel.ID,
		Servers: []*ServerToPost{
			{
				Server:   &server,
				Name:     name,
				SortName: getSortName(name),
			},
		},
		Category: category,
		SortBy:   sortBy,
	})
}

func (p *Plugin) getDiscordCategoryChannel(
	session *discord.Session,
	category *Category,
	name string,
) (*discordgo.Channel, error) {
	listGuild, err := p.state.Guild(category.GuildID)
	if err != nil {
		return nil, err
	}

	channelName := category.GroupBy.ChannelName(getSortName(name))

	for _, listGuildChannel := range listGuild.Channels {
		if listGuildChannel.ParentID != category.ChannelID {
			continue
		}

		if listGuildChannel.Name != channelName {
			continue
		}

		return listGuildChannel, nil
	}

	return session.Client.GuildChannelCreateComplex(
		category.GuildID,
		discordgo.GuildChannelCreateData{
			Name:     channelName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: category.ChannelID,
		},
	)
}

func (p *Plugin) getList(botID string) ([]*Server, error) {
	return serversFindMany(
		p.db,
		"state = ? AND bot_id = ?", StatePublic, botID,
	)
}

func channelServersToPostContain(key *ChannelServersToPost, list []*ChannelServersToPost) bool {
	for _, item := range list {
		if key.ChannelID != item.ChannelID {
			continue
		}

		return true
	}

	return false
}

func getCategoyTopic(category *Category, sortBy []SortBy) string {
	var sortByText string
	for _, item := range sortBy {
		sortByText += string(item) + ", "
	}
	sortByText = strings.TrimRight(sortByText, ", ")

	return fmt.Sprintf(
		"Keywords: %s\nSorted By: %s",
		strings.Join(category.Keywords, ", "),
		sortByText,
	)
}
