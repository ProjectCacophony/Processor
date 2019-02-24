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

	"github.com/golang/freetype"

	"golang.org/x/image/font"

	"golang.org/x/image/draw"
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

	fontType, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	freetypeContext := freetype.NewContext()
	freetypeContext.SetFont(fontType)
	freetypeContext.SetDPI(72)
	freetypeContext.SetHinting(font.HintingFull)
	freetypeContext.SetSrc(image.White)
	freetypeContext.SetDst(img)
	freetypeContext.SetClip(img.Bounds())

	fontOffsetX := (tileWidth / 100) * 2
	fontOffsetY := (tileHeight / 100) * 2
	var posX, posY int
	var tileRectangle image.Rectangle
	var fontSize float64
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
			fontSize = 36

			// TODO: measure text, and reduce size if needed
			// wait for https://godoc.org/golang.org/x/image/font/opentype ?
			// TODO: Border for better readability

			freetypeContext.SetFontSize(fontSize)
			freetypeContext.SetClip(tileRectangle.Bounds())
			pt := freetype.Pt(posX+fontOffsetX, posY+fontOffsetY+int(freetypeContext.PointToFixed(fontSize)>>6))

			freetypeContext.DrawString(descriptions[i], pt) // nolint: errcheck
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
