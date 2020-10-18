package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type enhancedEmoji struct {
	Emoji *discordgo.Emoji
	Link  string
}

var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

func (p *Plugin) handleDownloadEmoji(event *events.Event) {
	if !permissions.DiscordAttachFiles.Match(event.State(), event.DB(), event.BotUserID, event.ChannelID, false, false) {
		event.Respond("tools.download-emoji.no-attach-files-permission")
		return
	}

	messages, err := event.Respond("tools.download-emoji.preparing")
	if err != nil {
		event.Except(err)
		return
	}
	defer func() {
		for _, message := range messages {
			err := event.Discord().Client.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				break
			}
		}
	}()

	guild, err := event.State().Guild(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	emojiList := make([]enhancedEmoji, len(guild.Emojis))

	for i, emoji := range guild.Emojis {
		emojiList[i].Emoji = emoji
		emojiList[i].Link = fmt.Sprintf("%semojis/%s.png", discordgo.EndpointCDN, emoji.ID)

		if emoji.Animated {
			emojiList[i].Link = fmt.Sprintf("%semojis/%s.gif", discordgo.EndpointCDN, emoji.ID)
		}
	}

	if len(emojiList) <= 0 {
		event.Respond("tools.download-emoji.none-emoji")
		return
	}

	tempDir, err := ioutil.TempDir("", "emoji")
	if err != nil {
		event.Except(err)
		return
	}

	defer func() {
		err := os.RemoveAll(tempDir)
		event.ExceptSilent(err)
	}()

	zipName := fmt.Sprintf("cacophony-emoji-%s", alphanumericRegex.ReplaceAllString(guild.Name, ""))

	zipPath := filepath.Join(tempDir, fmt.Sprintf("%s.zip", zipName))

	outFile, err := os.Create(zipPath)
	if err != nil {
		event.Except(err)
		return
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	emojiMap := make(map[string]bool)

	var filename string
	var filenameI int
	for _, emoji := range emojiList {
		extension := filepath.Ext(emoji.Link)

		resp, err := event.HTTPClient().Get(emoji.Link)
		if err != nil {
			event.Except(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			event.Except(fmt.Errorf("expected status code 200, received %d", resp.StatusCode))
			return
		}

		filename = emoji.Emoji.Name + extension
		filenameI = 0
		for {
			if !emojiMap[filename] {
				break
			}
			filenameI++
			filename = emoji.Emoji.Name + strconv.Itoa(filenameI) + extension
		}

		emojiMap[filename] = true
		file, err := zipWriter.Create(filepath.Join(zipName, filename))
		if err != nil {
			event.Except(err)
			return
		}

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			event.Except(err)
			return
		}

		resp.Body.Close()
	}

	err = zipWriter.Close()
	if err != nil {
		event.Except(err)
		return
	}

	err = outFile.Close()
	if err != nil {
		event.Except(err)
		return
	}

	file, err := os.Open(zipPath)
	if err != nil {
		event.Except(err)
		return
	}
	defer file.Close()

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Content: "tools.download-emoji.success",
		Files: []*discordgo.File{
			{
				Name:   fmt.Sprintf("%s.zip", zipName),
				Reader: file,
			},
		},
	}, "userID", event.UserID)
	if err != nil {
		event.Except(err)
		return
	}
}
