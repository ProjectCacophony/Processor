package color

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleColor(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("color.too-few-arguments")
		return
	}

	re := regexp.MustCompile("[0-9]+")

	// remove # if needed
	hexText := strings.TrimPrefix(event.Fields()[1], "#")

	// if more than one field was passed to command, combine and attempt to parse out rgb values
	if len(event.Fields()) > 2 {

		var err error
		parsedNums := re.FindAllString(strings.Join(event.Fields(), ","), -1)
		hexText, err = getHexFromRGB(parsedNums)

		if err != nil {
			event.Respond("color.invalid-color")
			return
		}
	}

	// handle size 3 hex codes. eg. #345
	if len(hexText) == 3 {
		hexText = fmt.Sprintf("%c%c%c%c%c%c", hexText[0], hexText[0], hexText[1], hexText[1], hexText[2], hexText[2])
	}

	// decode hex to rbg
	rgbArray, err := hex.DecodeString(hexText)
	if err != nil {

		// if hex decode failed, attempt to parse rgb values from string. ex: 255,255,255
		hexText, _ = getHexFromRGB(re.FindAllString(event.Fields()[1], -1))

		rgbArray, err = hex.DecodeString(hexText)

		// if attempt failed, continue processing error from hex decode
		if err != nil {
			if strings.Contains(err.Error(), "odd length hex string") ||
				strings.Contains(err.Error(), "invalid byte") {
				event.Respond("color.invalid-color")
				return
			}

			event.Except(err)
			return
		}
	}

	if len(rgbArray) < 3 {
		event.Respond("color.invalid-color")
		return
	}

	// make square image with given color
	r := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(r)
	imageColor := color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{imageColor}, image.Point{}, draw.Src)
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

func getHexFromRGB(rgbValues []string) (string, error) {
	if len(rgbValues) != 3 {
		return "", errors.New("invalid rgb array size. expected size: 3")
	}

	validRGB := true
	hexArray := make([]string, 3)
	for _, rgbValue := range rgbValues {

		rgbNum, err := strconv.Atoi(rgbValue)
		if err != nil || rgbNum < 0 || rgbNum > 255 {
			validRGB = false
			break
		}

		hexValue := fmt.Sprintf("%X", rgbNum)
		if len(hexValue) == 1 {
			hexValue = "0" + hexValue
		}

		hexArray = append(hexArray, hexValue)
	}

	if !validRGB {
		return "", errors.New("rgb values not valid")
	}

	return strings.Join(hexArray, ""), nil
}
