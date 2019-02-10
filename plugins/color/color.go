package color

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func handleColor(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("color.too-few-arguments") // nolint: errcheck
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
			event.Respond("color.invalid-hex") // nolint: errcheck
			return
		}

		event.Except(err)
		return
	}

	if len(rgbArray) < 3 {
		event.Respond("color.invalid-hex") // nolint: errcheck
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
	err = jpeg.Encode(&buff, finalImage, nil)
	if err != nil {
		event.Except(err)
		return
	}

	// send image
	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: event.Translate("color.result",
				"hexcode", hexText,
				"r", rgbArray[0],
				"g", rgbArray[1],
				"b", rgbArray[2],
				"colorCode", discord.HexToColorCode(hexText)),
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + hexText + ".jpg",
			},
		},
		Files: []*discordgo.File{
			{
				Name:   hexText + ".jpg",
				Reader: bytes.NewReader(buff.Bytes()),
			},
		},
	})
	if err != nil {
		event.Except(err)
	}
}
