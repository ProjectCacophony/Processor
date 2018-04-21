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
		fmt.Printf("Error decoding hex string: %s \n", err.Error())
		return
	}

	// make square image with given color
	r := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(r)
	imageColor := color.RGBA{uint8(rgbArray[0]), uint8(rgbArray[1]), uint8(rgbArray[2]), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{imageColor}, image.ZP, draw.Src)
	finalImage := img.SubImage(r)

	// send image
	var buff bytes.Buffer
	png.Encode(&buff, finalImage)
	_, err = dhelpers.SendFile(event.MessageCreate.ChannelID, "test.png", bytes.NewReader(buff.Bytes()), "Color: #"+hexText)
	if err != nil {
		fmt.Println("Error sending file: ", err.Error())
	}
}
