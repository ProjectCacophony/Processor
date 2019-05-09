package serverlist

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCategoryCreate(event *events.Event) {
	if len(event.Fields()) < 6 {
		event.Respond("serverlist.too-few")
		return
	}

	fields := event.Fields()[3:]

	channel, fields, err := common.FieldsExtractChannelTypes(
		event, fields, discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildCategory,
	)
	if err != nil {
		event.Except(err)
		return
	}

	var keywords []string // nolint: prealloc
	for _, keyword := range strings.Split(fields[0], ";") {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			continue
		}

		keywords = append(keywords, keyword)
	}

	if len(keywords) == 0 {
		event.Respond("serverlist.category-create.no-keywords")
		return
	}

	var sortBys []SortBy
	for _, sortBy := range strings.Split(fields[1], ";") {
		sortBy = strings.ToLower(strings.TrimSpace(sortBy))

		for _, allSortBy := range allSortBys {
			if sortBy != string(allSortBy) {
				continue
			}

			sortBys = append(sortBys, allSortBy)
		}
	}

	if len(sortBys) == 0 {
		event.Respond("serverlist.category-create.invalid-sortby")
		return
	}

	var groupBy GroupBy
	if channel.Type == discordgo.ChannelTypeGuildCategory {

		for _, allGroupBy := range allGroupBys {
			if fields[2] != string(allGroupBy) {
				continue
			}

			groupBy = allGroupBy
		}

		if groupBy == "" {
			event.Respond("serverlist.category-create.invalid-groupby")
			return
		}

	}

	err = categoryCreate(
		p.db,
		keywords,
		event.BotUserID,
		event.GuildID,
		channel.ID,
		event.UserID,
		sortBys,
		groupBy,
	)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("serverlist.category-create.category-created")
	event.Except(err)
}
