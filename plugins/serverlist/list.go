package serverlist

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleListRefresh(event *events.Event) {
	err := p.refreshList(event.BotUserID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.list-refresh.content")
	event.Except(err)
}

func (p *Plugin) refreshList(botID string) error {
	allList, err := p.getList(botID)
	if err != nil {
		return err
	}

	session, err := discord.NewSession(p.tokens, botID)
	if err != nil {
		return errors.Wrap(err, "failure creating Discord Session for Bot")
	}

	var serversToPost []*ChannelServersToPost

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
						return err
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
			return err
		}
	}

	return nil
}

type ChannelServersToPost struct {
	ChannelID string
	Servers   []*ServerToPost
	SortBy    []SortBy
}

type ServerToPost struct {
	Server *Server
	Name   string
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
				return err
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
				return err
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
				return err
			}
		}
	}

	return p.setListMessages(channelToPost.ChannelID, messages)
}

func (p *Plugin) newListMessage(
	messages []*discordgo.Message,
	session *discord.Session,
	channelID string,
	content string,
) ([]*discordgo.Message, error) {
	message, err := discord.SendComplexWithVars(
		p.redis,
		session,
		nil,
		channelID,
		&discordgo.MessageSend{
			Content: content,
		},
		false,
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
	return fmt.Sprintf(
		"**%s** — https://discord.gg/%s\n%s",
		name,
		server.InviteCode,
		server.Description,
	)
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
			Server: &server,
			Name:   name,
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
				Server: &server,
				Name:   name,
			},
		},
		SortBy: sortBy,
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

	channelName := category.GroupBy.ChannelName(name)

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