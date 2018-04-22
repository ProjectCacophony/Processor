package stats

import (
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"gitlab.com/project-d-collab/SqsProcessor/metrics"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/state"
)

// displays various bot stats
func displayStats(event dhelpers.EventContainer) {
	botUser, err := state.User(event.BotUserID)
	dhelpers.CheckErr(err)

	allGuildIDs, err := state.AllGuildIDs()
	dhelpers.CheckErr(err)

	allUserIDs, err := state.AllUserIDs()
	dhelpers.CheckErr(err)

	var ram runtime.MemStats
	runtime.ReadMemStats(&ram)

	bootTime, err := strconv.ParseInt(metrics.Uptime.String(), 10, 64)
	if err != nil {
		bootTime = 0
	}

	statsEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username,
			IconURL: botUser.AvatarURL("64"),
		},
		Fields: []*discordgo.MessageEmbedField{},
	}

	// display gateway uptime
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "⌛ Gateway Uptime",
		Value:  dhelpers.HumanizeDuration(time.Since(event.GatewayStarted)),
		Inline: true,
	})

	// display sqs processor uptime
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "⌛ SqsP Uptime",
		Value:  dhelpers.HumanizeDuration(time.Since(time.Unix(bootTime, 0))),
		Inline: true,
	})

	// display sqs processor running coroutines
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "🔄 SqsP Coroutines",
		Value:  strconv.Itoa(runtime.NumGoroutine()),
		Inline: true,
	})

	// display sqs processor memory stats
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "💡 SqsP Memory",
		Value:  "Heap / Sys " + humanize.Bytes(ram.Alloc) + "/" + humanize.Bytes(ram.Sys),
		Inline: true,
	})

	// display sqs processor garbage collected
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "♻ SqsP GC",
		Value:  humanize.Bytes(ram.TotalAlloc),
		Inline: true,
	})

	// display sqs processor go version
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "♻ SqsP Go",
		Value:  runtime.Version(),
		Inline: true,
	})

	// display amount of guilds and channels
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name: "👥 Guilds / Users",
		Value: humanize.Comma(int64(len(allGuildIDs))) + "\n" +
			humanize.Comma(int64(len(allUserIDs))),
		Inline: true,
	})

	_, err = dhelpers.SendEmbed(event.MessageCreate.ChannelID, statsEmbed)
	dhelpers.CheckErr(err)
}
