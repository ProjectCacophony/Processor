package lastfm

// nolint: golint
import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Creates a Collage PNG Image from internet image urls (PNG or JPEG).
// imageUrls		: a slice with all image URLs. Empty strings will create an empty space in the collage.
// descriptions		: a slice with text that will be written on each tile. Can be empty.
// width			: the width of the result collage image.
// height			: the height of the result collage image.
// tileWidth		: the width of each tile image.
// tileHeight		: the height of each tile image.
// backgroundColour	: the background colour as a hex string.
func CollageFromURLs(
	client *http.Client, imageUrls, descriptions []string, width, height, tileWidth, tileHeight int,
) ([]byte, error) {
	imageDataArray := make([][]byte, 0)
	// download images
	for _, imageURL := range imageUrls {
		if imageURL == "" {
			imageDataArray = append(imageDataArray, nil)
			continue
		}

		res, err := client.Get(imageURL)
		if err != nil {
			imageDataArray = append(imageDataArray, nil)
			continue
		}

		byteData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			imageDataArray = append(imageDataArray, nil)
			continue
		}

		res.Body.Close()

		imageDataArray = append(imageDataArray, byteData)
	}

	return CollageFromBytes(imageDataArray, descriptions, width, height, tileWidth, tileHeight)
}

// Creates a Collage PNG Image from image []byte (PNG or JPEG).
// imageDataArray   : a slice of all image []byte data
// descriptions		: a slice with text that will be written on each tile. Can be empty.
// width			: the width of the result collage image.
// height			: the height of the result collage image.
// tileWidth		: the width of each tile image.
// tileHeight		: the height of each tile image.
// backgroundColour	: the background colour as a hex string.
func CollageFromBytes(
	imageDataArray [][]byte, descriptions []string, width, height, tileWidth, tileHeight int,
) ([]byte, error) {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	})
	backgroundColor := color.RGBA{R: 54, B: 57, G: 63, A: 0} // Discord Dark Theme background
	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.ZP, draw.Src)

	fontBytes, err := ioutil.ReadFile("assets/fonts/Spoqa Han Sans JP Bold.ttf")
	if err != nil {
		return nil, err
	}

	fontType, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	fontOffsetX := (tileWidth / 100) * 2
	fontOffsetY := (tileHeight / 100) * 2
	var posX, posY int
	var tileRectangle image.Rectangle
	var fontSize int
	var textPosX, textPosY int
	var textLength int
	for i, imageItemData := range imageDataArray {
		// switch tile to new line if required
		if posX > 0 && posX+tileWidth > width {
			posY += tileHeight
			posX = 0
		}
		tileRectangle = image.Rectangle{
			Min: image.Point{X: posX, Y: posY},
			Max: image.Point{X: posX + tileWidth, Y: posY + tileHeight},
		}

		imageItem, _, err := image.Decode(bytes.NewReader(imageItemData))
		if err == nil {

			draw.NearestNeighbor.Scale(img, tileRectangle, imageItem, imageItem.Bounds(), draw.Over, nil)
		}

		if len(descriptions) > i {
			for fontSize = (tileHeight / 100) * 15; fontSize > 1; fontSize-- {

				textPosX = posX + fontOffsetX
				textPosY = posY + fontOffsetY + fontSize - 5

				textLength, _ = measureString(
					fontType, descriptions[i],
					fontSize,
				)
				if textLength < tileWidth-fontOffsetY {
					break
				}
			}

			drawStringWithOutline(
				img,
				image.White, image.Black,
				fontType, descriptions[i],
				fontSize,
				textPosX, textPosY,
				5.0,
			)
		}

		posX += tileWidth
	}

	var buffer bytes.Buffer

	err = jpeg.Encode(io.Writer(&buffer), img, &jpeg.Options{
		Quality: 95,
	})
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func drawStringWithOutline(
	dst draw.Image,
	outlineSrc image.Image,
	textSrc image.Image,
	fontType *truetype.Font,
	text string,
	fontSize int,
	x, y int,
	outlineSize int,
) {
	for dy := -outlineSize; dy <= outlineSize; dy++ {
		for dx := -outlineSize; dx <= outlineSize; dx++ {
			if dx*dx+dy*dy >= outlineSize*outlineSize {
				// give it rounded corners
				continue
			}

			drawString(
				dst, outlineSrc,
				fontType, text,
				fontSize,
				x+dx, y+dy,
			)
		}
	}

	drawString(
		dst, textSrc,
		fontType, text,
		fontSize,
		x, y,
	)
}

func drawString(
	dst draw.Image,
	src image.Image,
	fontType *truetype.Font,
	text string,
	fontSize int,
	x, y int,
) {
	drawer := &font.Drawer{
		Dst: dst,
		Src: src,
		Face: truetype.NewFace(fontType, &truetype.Options{
			Size:    float64(fontSize),
			Hinting: font.HintingFull,
		}),
		Dot: fixed.P(x, y),
	}

	drawer.DrawString(text)
}

func measureString(
	fontType *truetype.Font,
	text string,
	fontSize int,
) (int, int) {
	drawer := &font.Drawer{
		Face: truetype.NewFace(fontType, &truetype.Options{
			Size:    float64(fontSize),
			Hinting: font.HintingFull,
		}),
	}

	res := drawer.MeasureString(text)
	return res.Round(), fontSize
}
