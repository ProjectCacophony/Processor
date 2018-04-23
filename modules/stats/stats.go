package stats

import (
	"runtime"
	"strconv"
	"time"

	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"gitlab.com/project-d-collab/SqsProcessor/metrics"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/cache"
	"gitlab.com/project-d-collab/dhelpers/state"
)

// displays various bot stats
func displayStats(event dhelpers.EventContainer) {
	// read information from stats
	botUser, err := state.User(event.BotUserID)
	dhelpers.CheckErr(err)

	allGuildIDs, err := state.AllGuildIDs()
	dhelpers.CheckErr(err)

	allUserIDs, err := state.AllUserIDs()
	dhelpers.CheckErr(err)

	// read information from the runtime
	var ram runtime.MemStats
	runtime.ReadMemStats(&ram)

	// read information from metrics
	bootTime, err := strconv.ParseInt(metrics.Uptime.String(), 10, 64)
	if err != nil {
		bootTime = 0
	}

	// read information from redis
	var redisConnectedClients, redisUsedMemoryHuman string

	redisInfoText, err := cache.GetRedisClient().Info().Result()
	dhelpers.CheckErr(err)
	for _, redisInfoLine := range strings.Split(redisInfoText, "\r\n") {
		args := strings.Split(redisInfoLine, ":")
		if len(args) < 2 {
			continue
		}

		switch args[0] {
		case "connected_clients":
			redisConnectedClients = args[1]
		case "used_memory_human":
			redisUsedMemoryHuman = args[1]
		}
	}

	// build embed
	statsEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username + " Statistics",
			IconURL: botUser.AvatarURL("64"),
		},
		Fields: []*discordgo.MessageEmbedField{},
	}

	// display amount of guilds and channels
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name: "👥 Guilds / Users",
		Value: humanize.Comma(int64(len(allGuildIDs))) + "\n" +
			humanize.Comma(int64(len(allUserIDs))),
		Inline: true,
	})

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
		Value:  "Heap " + humanize.Bytes(ram.Alloc) + "\nSys " + humanize.Bytes(ram.Sys),
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
		Name:   "📌 SqsP Go",
		Value:  strings.Replace(runtime.Version(), "go", "v", 1),
		Inline: true,
	})

	// display redis information
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "📋 Redis",
		Value:  "Clients " + redisConnectedClients + "\nMemory " + redisUsedMemoryHuman,
		Inline: true,
	})

	// send it
	_, err = dhelpers.SendEmbed(event.MessageCreate.ChannelID, statsEmbed)
	dhelpers.CheckErr(err)
}
