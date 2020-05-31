package quickactions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
)

const questionnaireRemindKey = "cacophony:processor:quickactions:remind"

var remindEmoji = map[string]time.Duration{
	"quickaction_remind_custom": 0, // custom delay
	"quickaction_remind_1h":     1 * time.Hour,
	"quickaction_remind_8h":     8 * time.Hour,
	"quickaction_remind_24h":    24 * time.Hour,
}

func (p *Plugin) remindMessage(event *events.Event) {
	params := quickactionParams{
		GuildID:   event.MessageReactionAdd.GuildID,
		ChannelID: event.MessageReactionAdd.ChannelID,
		MessageID: event.MessageReactionAdd.MessageID,
		Emoji:     &event.MessageReactionAdd.Emoji,
		ToUserID:  event.MessageReactionAdd.UserID,
		BotUserID: event.BotUserID,
		Delay:     remindEmoji[event.MessageReactionAdd.Emoji.Name],
	}
	if params.Delay <= 0 {
		p.remindAskCustomDelay(event, params)
		return
	}

	err := p.setupQuickactionRemind(
		event.Context(),
		event.State(),
		event.Discord(),
		params,
	)
	if err != nil {
		event.ExceptSilent(err)
	}
}

func (p *Plugin) remindAskCustomDelay(
	event *events.Event,
	params quickactionParams,
) {
	data, err := params.Marshal()
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	messages, err := event.Send(
		event.ChannelID,
		"quickactions.remind.ask-custom-delay",
		"UserID",
		event.MessageReactionAdd.UserID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
	if len(messages) <= 0 {
		return
	}

	err = event.Questionnaire().Register(
		questionnaireRemindKey,
		events.QuestionnaireFilter{
			GuildID:   event.GuildID,
			ChannelID: event.ChannelID,
			UserID:    event.MessageReactionAdd.UserID,
			Type:      events.MessageCreateType,
		},
		map[string]interface{}{
			"params":                      string(data),
			"question_message_channel_id": messages[0].ChannelID,
			"question_message_id":         messages[0].ID,
		},
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
}

func (p *Plugin) handleRemindQuestionnaire(event *events.Event) {
	var params quickactionParams
	err := params.Unmarshal([]byte(event.QuestionnaireMatch.Payload["params"].(string)))
	if err != nil {
		event.ExceptSilent(err)
	}

	duration, err := time.ParseDuration(event.MessageCreate.Content)
	if err != nil || duration < 1*time.Minute {
		messages, _ := event.Send(
			event.ChannelID,
			"quickactions.remind.invalid-custom-delay",
		)

		go func() {
			time.Sleep(5 * time.Second)

			// delete questionnaire response message if possible
			if permissions.DiscordManageMessages.Match(
				event.State(),
				p.db,
				params.BotUserID,
				params.ChannelID,
				false,
				false,
			) {
				event.Discord().Client.ChannelMessageDelete(
					event.ChannelID,
					event.MessageID,
				)
			}

			// delete questionnaire question message
			event.Discord().Client.ChannelMessageDelete(
				event.QuestionnaireMatch.Payload["question_message_channel_id"].(string),
				event.QuestionnaireMatch.Payload["question_message_id"].(string),
			)

			// delete questionnaire error message
			for _, message := range messages {
				event.Discord().Client.ChannelMessageDelete(
					message.ChannelID,
					message.ID,
				)
			}

			// remove reaction if possible
			if permissions.DiscordManageMessages.Match(
				event.State(),
				p.db,
				params.BotUserID,
				params.ChannelID,
				false,
				false,
			) {
				event.Discord().Client.MessageReactionRemove(
					params.ChannelID,
					params.MessageID,
					params.Emoji.APIName(),
					params.ToUserID,
				)
			}
		}()
		return
	}

	params.Delay = duration

	err = p.setupQuickactionRemind(
		event.Context(),
		event.State(),
		event.Discord(),
		params,
	)
	if err != nil {
		event.ExceptSilent(err)
	}

	go func() {
		time.Sleep(1 * time.Second)

		// delete questionnaire response message if possible
		if permissions.DiscordManageMessages.Match(
			event.State(),
			p.db,
			params.BotUserID,
			params.ChannelID,
			false,
			false,
		) {
			event.Discord().Client.ChannelMessageDelete(
				event.ChannelID,
				event.MessageID,
			)
		}

		// delete questionnaire question message
		event.Discord().Client.ChannelMessageDelete(
			event.QuestionnaireMatch.Payload["question_message_channel_id"].(string),
			event.QuestionnaireMatch.Payload["question_message_id"].(string),
		)
	}()
}

type quickactionParams struct {
	GuildID   string
	ChannelID string
	MessageID string
	Emoji     *discordgo.Emoji
	ToUserID  string
	BotUserID string
	Delay     time.Duration
}

func (qp *quickactionParams) Marshal() ([]byte, error) {
	return json.Marshal(qp)
}

func (qp *quickactionParams) Unmarshal(data []byte) error {
	return json.Unmarshal(data, qp)
}

func (p *Plugin) setupQuickactionRemind(
	ctx context.Context,
	state *state.State,
	session *discord.Session,
	params quickactionParams,
) error {
	newEvent, err := events.New(events.CacophonyQuickactionRemind)
	if err != nil {
		return err
	}
	newEvent.QuickactionRemind = &events.QuickactionRemind{
		GuildID:   params.GuildID,
		ChannelID: params.ChannelID,
		MessageID: params.MessageID,
		Emoji:     params.Emoji,
		ToUserID:  params.ToUserID,
	}
	newEvent.BotUserID = params.BotUserID

	err = p.publisher.PublishAt(
		ctx,
		newEvent,
		time.Now().Add(params.Delay),
	)
	if err != nil {
		return err
	}

	ackEmoji := emoji.GetWithout("ok")

	err = discord.React(
		p.redis,
		session,
		params.ChannelID,
		params.MessageID,
		false,
		ackEmoji,
	)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(3 * time.Second)

		// remove reaction if possible
		if permissions.DiscordManageMessages.Match(
			state,
			p.db,
			params.BotUserID,
			params.ChannelID,
			false,
			false,
		) {
			session.Client.MessageReactionRemove(
				params.ChannelID,
				params.MessageID,
				params.Emoji.APIName(),
				params.ToUserID,
			)
		}

		session.Client.MessageReactionRemove(
			params.ChannelID,
			params.MessageID,
			ackEmoji,
			params.BotUserID,
		)
	}()

	return nil
}

func (p *Plugin) handleQuickactionRemind(event *events.Event) {
	message, err := discord.FindMessage(
		event.State(),
		event.Discord(),
		event.QuickactionRemind.ChannelID,
		event.QuickactionRemind.MessageID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	message.GuildID = event.QuickactionRemind.GuildID // :(

	// copy message into DMs
	_, err = event.SendComplexDM(
		event.QuickactionRemind.ToUserID,
		&discordgo.MessageSend{
			Content: "quickactions.message.content",
			Embed:   convertMessageToEmbed(message),
		},
		"message",
		message,
		"emoji",
		event.QuickactionRemind.Emoji.MessageFormat(),
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
}
