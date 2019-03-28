package serverlist

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/permissions"

	humanize "github.com/dustin/go-humanize"

	lock "github.com/bsm/redis-lock"

	"gitlab.com/Cacophony/go-kit/discord"

	"github.com/bwmarrin/discordgo"

	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) handleQueue(event *events.Event) {
	err := config.GuildSetString(
		p.db, event.GuildID, queueChannelIDKey, event.ChannelID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue.content", "channelID", event.ChannelID)
	event.Except(err)

	p.refreshQueue(event.GuildID)
}

func (p *Plugin) handleQueueReaction(event *events.Event) bool {
	if event.MessageReactionAdd.Emoji.Name != emojiApprove {
		return false
	}

	channelID, err := config.GuildGetString(
		p.db, event.GuildID, queueChannelIDKey,
	)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.ExceptSilent(err)
		}
		return false
	}

	if event.ChannelID != channelID {
		return false
	}

	if !event.Has(permissions.BotOwner) {
		return false
	}

	err = p.approveCurrentServer(event.BotUserID, event.GuildID)
	if err != nil &&
		!strings.Contains(err.Error(), "nothing to approve") &&
		!strings.Contains(err.Error(), "nothing to reject") {
		event.ExceptSilent(err)
		return true
	}

	p.refreshQueue(event.GuildID)

	discord.RemoveReact( // nolint: errcheck
		p.redis,
		event.Discord(),
		event.ChannelID,
		event.MessageReactionAdd.MessageID,
		event.UserID,
		false,
		event.MessageReactionAdd.Emoji.Name,
	)

	return true
}

func (p *Plugin) handleQueueReject(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.queue-reject.no-reason") // nolint: errcheck
		return
	}

	err := p.rejectCurrentServer(event.BotUserID, event.GuildID, event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue-reject.success")
	event.Except(err)
}

func (p *Plugin) handleQueueRefresh(event *events.Event) {
	err := p.refreshQueueForGuild(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.queue-refresh.content")
	event.Except(err)
}

func (p *Plugin) refreshQueue(guildIDs ...string) {
	var err error

	for _, guildID := range guildIDs {
		err = p.refreshQueueForGuild(guildID)
		if err != nil {
			p.logger.Error("failure refreshing queue for guild",
				zap.Error(err),
				zap.String("guild_id", guildID),
			)
		}
	}
}

func (p *Plugin) approveCurrentServer(botID, guildID string) error {
	queueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	server := queueFind(queueMessage.CurrentServerID, queue)
	if server == nil {
		return errors.New("found nothing to approve")
	}

	return server.QueueApprove(p)
}

func (p *Plugin) rejectCurrentServer(botID, guildID, reason string) error {
	queueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	server := queueFind(queueMessage.CurrentServerID, queue)
	if server == nil {
		return errors.New("found nothing to reject")
	}

	return server.QueueReject(p, reason)
}

// TODO: refactor!
// nolint: gocyclo
func (p *Plugin) refreshQueueForGuild(guildID string) error {
	guildLock := p.getGuildLock(guildID)

	locked, err := guildLock.Lock()
	if err != nil {
		return errors.Wrap(err, "error acquiring lock")
	}
	if !locked {
		return errors.Wrap(err, "unable to acquire lock")
	}
	defer guildLock.Unlock() // nolint: errcheck

	channelID, err := config.GuildGetString(
		p.db, guildID, queueChannelIDKey,
	)
	if err != nil {
		return errors.Wrap(err, "failure getting queueChannelIDKey from guild config")
	}

	if channelID == "" {
		return errors.New("no queue channel ID set for guild")
	}

	botID, err := p.state.BotForGuild(guildID)
	if err != nil {
		return errors.Wrap(err, "failure getting Bot ID for Guild")
	}

	queue, err := p.getQueue(botID)
	if err != nil {
		return errors.Wrap(err, "unable to query for queued entries")
	}

	session, err := discord.NewSession(p.tokens, botID)
	if err != nil {
		return errors.Wrap(err, "failure creating Discord Session for Bot")
	}

	queueMessage, err := p.getCurrentQueueMessage(guildID)
	if err != nil {
		return errors.Wrap(err, "error getting QueueMessage from config")
	}

	if queueMessage == nil {
		if len(queue) == 0 {
			return nil
		}

		embed := getQueueMessageEmbed(queue[0], len(queue))

		messages, err := discord.SendComplexWithVars(
			p.redis,
			session,
			p.Localisations(),
			channelID,
			&discordgo.MessageSend{
				Embed: embed,
			},
			false,
		)
		if err != nil {
			return errors.Wrap(err, "error sending initial QueueMessage")
		}
		queueMessage = &QueueMessage{
			CurrentServerID: queue[0].ID,
			MessageID:       messages[0].ID,
			Embed:           embed,
		}

		err = config.GuildSetInterface(p.db, guildID, queueMessageKey, queueMessage)
		if err != nil {
			return errors.Wrap(err, "error saving initial QueueMessage")
		}

		discord.React( // nolint: errcheck
			p.redis,
			session,
			channelID,
			messages[0].ID,
			false,
			emojiApprove,
		)
	} else {
		item := queueFind(queueMessage.CurrentServerID, queue)
		if item != nil {
			embed := getQueueMessageEmbed(item, len(queue))

			_, err = discord.EditComplexWithVars(
				p.redis,
				session,
				p.Localisations(),
				&discordgo.MessageEdit{
					Embed:   embed,
					ID:      queueMessage.MessageID,
					Channel: channelID,
				},
				false,
			)
			if err != nil {
				if errD, ok := err.(*discordgo.RESTError); ok &&
					errD.Message != nil &&
					errD.Message.Code == discordgo.ErrCodeUnknownMessage {
					err = config.GuildSetInterface(p.db, guildID, queueMessageKey, nil)
					if err != nil {
						return errors.Wrap(err, "error saving empty")
					}
					return nil
				}

				return errors.Wrap(err, "error updating initial QueueMessage")
			}

			queueMessage.Embed = embed

			err = config.GuildSetInterface(p.db, guildID, queueMessageKey, queueMessage)
			if err != nil {
				return errors.Wrap(err, "error saving updated initial QueueMessage")
			}
		} else {
			if len(queue) > 0 {
				item = queue[0]
				embed := getQueueMessageEmbed(item, len(queue))

				_, err = discord.EditComplexWithVars(
					p.redis,
					session,
					p.Localisations(),
					&discordgo.MessageEdit{
						Embed:   embed,
						ID:      queueMessage.MessageID,
						Channel: channelID,
					},
					false,
				)
				if err != nil {
					if errD, ok := err.(*discordgo.RESTError); ok &&
						errD.Message != nil &&
						errD.Message.Code == discordgo.ErrCodeUnknownMessage {
						err = config.GuildSetInterface(p.db, guildID, queueMessageKey, nil)
						if err != nil {
							return errors.Wrap(err, "error saving empty")
						}
						return nil
					}

					return errors.Wrap(err, "error updating to new QueueMessage")
				}

				queueMessage.Embed = embed
				queueMessage.CurrentServerID = item.ID

				err = config.GuildSetInterface(p.db, guildID, queueMessageKey, queueMessage)
				if err != nil {
					return errors.Wrap(err, "error saving updated to new QueueMessage")
				}
			} else {
				embed := getQueueMessageEmbed(item, len(queue))

				_, err = discord.EditComplexWithVars(
					p.redis,
					session,
					p.Localisations(),
					&discordgo.MessageEdit{
						Embed:   embed,
						ID:      queueMessage.MessageID,
						Channel: channelID,
					},
					false,
				)
				if err != nil {
					if errD, ok := err.(*discordgo.RESTError); ok &&
						errD.Message != nil &&
						errD.Message.Code == discordgo.ErrCodeUnknownMessage {
						err = config.GuildSetInterface(p.db, guildID, queueMessageKey, nil)
						if err != nil {
							return errors.Wrap(err, "error saving empty")
						}
						return nil
					}

					return errors.Wrap(err, "error updating to new QueueMessage")
				}

				queueMessage.Embed = embed
				queueMessage.CurrentServerID = 0

				err = config.GuildSetInterface(p.db, guildID, queueMessageKey, queueMessage)
				if err != nil {
					return errors.Wrap(err, "error saving updated to new QueueMessage")
				}
			}
		}
	}

	return nil
}

func queueFind(n uint, list []*Server) *Server {
	for _, item := range list {
		if item.ID != n {
			continue
		}

		return item
	}

	return nil
}

func getQueueMessageEmbed(server *Server, total int) *discordgo.MessageEmbed {
	if server == nil {
		return &discordgo.MessageEmbed{
			Title:       "⌛ Serverlist Queue",
			Description: "Queue empty!",
		}
	}

	var categoryText string
	for _, category := range server.Categories {
		categoryText += "<#" + category.Category.ChannelID + ">, "
	}
	categoryText = strings.TrimRight(categoryText, ", ")

	return &discordgo.MessageEmbed{
		Title:       "⌛ Serverlist Queue",
		Description: "serverlist.queue.embed.description",
		Timestamp:   server.CreatedAt.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"there are %d Servers queued in total • added", total,
			),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "🏷 Name(s)",
				Value: fmt.Sprintf("%s\n#%s",
					strings.Join(server.Names, "; "), server.GuildID,
				),
				Inline: true,
			},
			{
				Name: "👥 Editor(s)",
				Value: fmt.Sprintf("<@%s>",
					strings.Join(server.EditorUserIDs, "> <@"),
				),
				Inline: true,
			},
			{
				Name:   "🚩 Invite",
				Value:  fmt.Sprintf("discord.gg/%s", server.InviteCode),
				Inline: true,
			},
			{
				Name:   "📈 Members",
				Value:  humanize.Comma(int64(server.TotalMembers)),
				Inline: true,
			},
			{
				Name:   "📖 Description",
				Value:  server.Description,
				Inline: false,
			},
			{
				Name:   "🗃 Category",
				Value:  categoryText,
				Inline: false,
			},
		},
	}
}

func (p *Plugin) getGuildLock(guildID string) *lock.Locker {
	return lock.New(
		p.redis,
		refreshQueueLock(guildID),
		&lock.Options{
			LockTimeout: 5 * time.Minute,
			RetryCount:  9, // try for 1 1/2 minutes
			RetryDelay:  10 * time.Second,
		},
	)
}

func (p *Plugin) getCurrentQueueMessage(guildID string) (*QueueMessage, error) {
	var queueMessage *QueueMessage
	err := config.GuildGetInterface(p.db, guildID, queueMessageKey, &queueMessage)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}

	return queueMessage, nil
}

func (p *Plugin) getQueue(botID string) ([]*Server, error) {
	return serversFindMany(
		p.db,
		"state = ? AND bot_id = ?", StateQueued, botID,
	)
}
