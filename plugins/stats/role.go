package stats

import (
	"bytes"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func findRole(event *events.Event) (*discordgo.Role, string, error) {
	if event.Has(permissions.BotAdmin) {
		for _, fieldA := range event.Fields() {
			for _, fieldB := range event.Fields() {

				role, err := event.State().RoleFromMention(fieldA, fieldB)
				if err == nil {
					return role, fieldA, nil
				}
			}
		}
	}

	role, err := event.FindRole()
	return role, event.GuildID, err
}

func (p *Plugin) handleRole(event *events.Event) {
	role, guildID, err := findRole(event)
	if err != nil {
		if strings.Contains(err.Error(), "role not found") {
			event.Respond("stats.role.not-found")
			return
		}
		event.Except(err)
		return
	}

	createdAt, err := discordgo.SnowflakeTimestamp(role.ID)
	if err != nil {
		event.Except(err)
		return
	}

	guild, err := event.State().Guild(guildID)
	if err != nil {
		event.Except(err)
		return
	}

	var colorCode string
	if role.Color != 0 {
		colorCode = discord.ColorCodeToHex(role.Color)
	}

	send := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "stats.role.embed.title",
			Description: "stats.role.embed.description",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "stats.role.embed.field.created-at.name",
					Value: "stats.role.embed.field.created-at.value",
				},
			},
			Color: role.Color,
		},
	}

	if colorCode != "" {
		image, err := generateColorImage(colorCode)
		if err == nil {
			send.Files = append(send.Files, &discordgo.File{
				Name:   colorCode + ".jpg",
				Reader: bytes.NewReader(image.Bytes()),
			})

			send.Embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: "attachment://" + colorCode + ".jpg",
			}
		}
	}

	_, err = event.RespondComplex(
		send,
		"role", role,
		"createdAt", createdAt,
		"guild", guild,
		"colorCode", colorCode,
	)
	event.Except(err)
}

func generateColorImage(hexCode string) (*bytes.Buffer, error) {
	// decode hex to rbg
	rgbArray, err := hex.DecodeString(hexCode)
	if err != nil {
		return nil, err
	}

	if len(rgbArray) < 3 {
		return nil, err
	}

	// make square image with given color
	r := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(r)
	imageColor := color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{imageColor}, image.ZP, draw.Src)
	finalImage := img.SubImage(r)

	// encode image
	var buff bytes.Buffer
	err = jpeg.Encode(&buff, finalImage, nil)
	if err != nil {
		return nil, err
	}

	return &buff, nil
}
