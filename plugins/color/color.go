package color

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func handleColor(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Send(event.MessageCreate.ChannelID, "color.too-few-arguments") // nolint: errcheck
		return
	}

	// remove # if needed
	hexText := strings.TrimPrefix(event.Fields()[1], "#")

	// handle size 3 hex codes. eg. #345
	if len(hexText) == 3 {
		hexText = fmt.Sprintf("%c%c%c%c%c%c", hexText[0], hexText[0], hexText[1], hexText[1], hexText[2], hexText[2])
	}

	// decode hex to rbg
	rgbArray, err := hex.DecodeString(hexText)
	if err != nil {
		if strings.Contains(err.Error(), "odd length hex string") ||
			strings.Contains(err.Error(), "invalid byte") {
			event.Send(event.MessageCreate.ChannelID, "color.invalid-hex") // nolint: errcheck
			return
		}

		event.Except(err)
		return
	}

	if len(rgbArray) < 3 {
		event.Send(event.MessageCreate.ChannelID, "color.invalid-hex") // nolint: errcheck
		return
	}

	// make square image with given color
	r := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(r)
	imageColor := color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{imageColor}, image.ZP, draw.Src)
	finalImage := img.SubImage(r)

	// encode image
	var buff bytes.Buffer
	err = png.Encode(&buff, finalImage)
	if err != nil {
		event.Except(err)
		return
	}

	// send image
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: event.Translate("color.result", "hexcode", hexText),
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + hexText + ".png",
			},
		},
		Files: []*discordgo.File{
			{
				Name:   hexText + ".png",
				Reader: bytes.NewReader(buff.Bytes()),
			},
		},
	})
	if err != nil {
		event.Except(err)
	}
}
