package stats

import (
	"strconv"
	"time"

	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/apihelper"
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

	allChannelIDs, err := state.AllChannelIDs()
	dhelpers.CheckErr(err)

	allUserIDs, err := state.AllUserIDs()
	dhelpers.CheckErr(err)

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

	// read sqsProcessor information
	var sqsProcessorText string
	sqsProcessorStatuses := apihelper.ReadSqsProcessorStatus()
	sqsProcessorText += "(" + strconv.Itoa(len(sqsProcessorStatuses)) + ")\n"
	for _, sqsProcessorStatusesEntry := range sqsProcessorStatuses {
		if !sqsProcessorStatusesEntry.Available {
			sqsProcessorText += "dead\n" // nolint: goconst
			if !strings.HasPrefix(sqsProcessorText, "⚠ ") {
				sqsProcessorText = "⚠ " + sqsProcessorText
			}
			continue
		}
		sqsProcessorText += "CR " + strconv.Itoa(sqsProcessorStatusesEntry.Service.Coroutines) + " "
		sqsProcessorText += "Mem " + humanize.Bytes(sqsProcessorStatusesEntry.Service.Heap) + "\n"
		sqsProcessorText += "Uptime " + dhelpers.HumanizeDuration(time.Since(sqsProcessorStatusesEntry.Service.Launch)) + "\n"
	}

	// read worker information
	var workerText string
	workersStatuses := apihelper.ReadWorkerStatus()
	workerText += "(" + strconv.Itoa(len(workersStatuses)) + ")\n"
	for _, workersStatusesEntry := range workersStatuses {
		if !workersStatusesEntry.Available {
			workerText += "dead\n" // nolint: goconst
			if !strings.HasPrefix(workerText, "⚠ ") {
				workerText = "⚠ " + workerText
			}
			continue
		}
		workerText += "Jobs " + strconv.Itoa(len(workersStatusesEntry.Entries)) + " "
		workerText += "CR " + strconv.Itoa(workersStatusesEntry.Service.Coroutines) + " "
		workerText += "Mem " + humanize.Bytes(workersStatusesEntry.Service.Heap) + "\n"
		workerText += "Uptime " + dhelpers.HumanizeDuration(time.Since(workersStatusesEntry.Service.Launch)) + "\n"
	}
	workerText = strings.TrimRight(workerText, "\n")

	// read gateway information
	var gatewayText string
	gatewayStatuses := apihelper.ReadGatewayStatus()
	gatewayText += "(" + strconv.Itoa(len(gatewayStatuses)) + ")\n"
	for _, gatewayStatusEntry := range gatewayStatuses {
		if !gatewayStatusEntry.Available {
			gatewayText += "dead\n" // nolint: goconst
			if !strings.HasPrefix(gatewayText, "⚠ ") {
				gatewayText = "⚠ " + gatewayText
			}
			continue
		}
		gatewayText += "CR " + strconv.Itoa(gatewayStatusEntry.Service.Coroutines) + " "
		gatewayText += "Mem " + humanize.Bytes(gatewayStatusEntry.Service.Heap) + "\n"
		gatewayText += "Uptime " + dhelpers.HumanizeDuration(time.Since(gatewayStatusEntry.Service.Launch)) + "\n"
	}
	gatewayText = strings.TrimRight(gatewayText, "\n")

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
		Name: "👥 State",
		Value: "G " + humanize.Comma(int64(len(allGuildIDs))) + "\n" +
			"C " + humanize.Comma(int64(len(allChannelIDs))) + "\n" +
			"U " + humanize.Comma(int64(len(allUserIDs))),
		Inline: true,
	})

	// display worker information
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "📠 SqsProcessor",
		Value:  sqsProcessorText,
		Inline: true,
	})

	// display worker information
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "👷 Worker",
		Value:  workerText,
		Inline: true,
	})

	// display gateway information
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "🚪 Gateway",
		Value:  gatewayText,
		Inline: true,
	})

	// display redis information
	statsEmbed.Fields = append(statsEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   "📋 Redis",
		Value:  "Clients " + redisConnectedClients + "\nMemory " + redisUsedMemoryHuman,
		Inline: true,
	})

	// send it
	_, err = event.SendEmbed(event.MessageCreate.ChannelID, statsEmbed)
	dhelpers.CheckErr(err)
}
