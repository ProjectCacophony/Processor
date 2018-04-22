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
	"gitlab.com/project-d-collab/dhelpers"
)

// takes hex code and displays matching color
func displayColor(event dhelpers.EventContainer) {
	if len(event.Args) < 2 {
		return
	}

	// remove # if needed
	hexText := strings.TrimPrefix(event.Args[1], "#")

	// handle size 3 hex codes. eg. #345
	if len(hexText) == 3 {
		hexText = fmt.Sprintf("%c%c%c%c%c%c", hexText[0], hexText[0], hexText[1], hexText[1], hexText[2], hexText[2])
	}

	// decode hex to rbg
	rgbArray, err := hex.DecodeString(hexText)
	if err != nil {
		if strings.Contains(err.Error(), "odd length hex string") ||
			strings.Contains(err.Error(), "invalid byte") {
			dhelpers.SendMessage(event.MessageCreate.ChannelID, "ColorInvalidHex") // nolint: errcheck, gas
			return
		}
	}
	dhelpers.CheckErr(err)

	if len(rgbArray) < 3 {
		dhelpers.SendMessage(event.MessageCreate.ChannelID, "ColorInvalidHex") // nolint: errcheck, gas
		return
	}

	// make square image with given color
	r := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(r)
	imageColor := color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{imageColor}, image.ZP, draw.Src)
	finalImage := img.SubImage(r)

	// send image
	var buff bytes.Buffer
	err = png.Encode(&buff, finalImage)
	dhelpers.CheckErr(err)

	_, err = dhelpers.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: dhelpers.Tf("ColorResult", "hexcode", hexText),
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
	dhelpers.CheckErr(err)
}
