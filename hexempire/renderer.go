package hexempire

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const fileDir = "./images"

// portNeighborSprite indexes, per hex direction, which port sprite variant
// and orientation to draw so the port always faces the adjacent water.
var (
	portImageNum = [6]int{2, 1, 2, 2, 1, 2}
	portFlipH    = [6]int{1, 0, 0, 0, 0, 1}
	portFlipV    = [6]int{1, 1, 1, 0, 0, 0}
)

// MapRenderer draws a generated map's board data and precomputed sprite
// choices onto a HexMap's image layers. It never generates randomness of
// its own; everything it draws is already decided by MapGenerator.
type MapRenderer struct{}

func NewMapRenderer() *MapRenderer {
	return &MapRenderer{}
}

func (r *MapRenderer) Render(hexMap *HexMap) {
	board := hexMap.Board

	r.drawMapBackground(hexMap, board)
	r.drawWaterAndPorts(hexMap, board)
	r.drawTownBackground(hexMap, board)
	r.drawTowns(hexMap, board)
	r.drawTownNames(hexMap, board)
	r.composeBackground(hexMap, board)
	r.drawUI(hexMap)
}

func (r *MapRenderer) drawMapBackground(hexMap *HexMap, board *Board) {
	dirtBgImg := loadImages(fileDir+"/ld_%d.png", 6)
	grassBgImg := loadImages(fileDir+"/l_%d.png", 6)

	for x := 0; x < 6; x++ {
		for y := 0; y < 4; y++ {
			imageOptions := hexMap.backgroundImageOptions[fieldKey(x, y)]
			dirtImg := dirtBgImg[imageOptions.BgDirtImgIndex]
			grassImg := grassBgImg[imageOptions.BgGrassImgIndex]

			destX := x*125 - 15
			destY := y*125 - 15

			// Dirt and grass each need their own transform: sharing one
			// DrawImageOptions and applying flip+rotate for both images in
			// sequence compounds the two transforms instead of applying
			// each independently (flip cancels out, rotation doubles).
			dirtOptions := &ebiten.DrawImageOptions{}
			flipImageMatrix(dirtOptions, dirtImg, imageOptions.FlipH, imageOptions.FlipV)
			rotateImageMatrix(dirtOptions, dirtImg, imageOptions.RotateDegrees)
			dirtOptions.GeoM.Translate(float64(destX), float64(destY))
			hexMap.BackgroundDirt.DrawImage(dirtImg, dirtOptions)

			grassOptions := &ebiten.DrawImageOptions{}
			flipImageMatrix(grassOptions, grassImg, imageOptions.FlipH, imageOptions.FlipV)
			rotateImageMatrix(grassOptions, grassImg, imageOptions.RotateDegrees)
			grassOptions.GeoM.Translate(float64(destX), float64(destY))
			hexMap.BackgroundGrass.DrawImage(grassImg, grassOptions)
		}
	}
}

func (r *MapRenderer) drawWaterAndPorts(hexMap *HexMap, board *Board) {
	waterBgImg := loadImages(fileDir+"/m_%d.png", 6)
	portBgImg := loadImages(fileDir+"/m_p%d.png", 2)

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Type != FieldTypeWater {
			return
		}

		imageOptions := hexMap.waterImageOptions[fieldKey(x, y)]
		waterImg := waterBgImg[imageOptions.WaterBgImgIndex]

		options := &ebiten.DrawImageOptions{}
		flipImageMatrix(options, waterImg, imageOptions.FlipH, imageOptions.FlipV)
		rotateImageMatrix(options, waterImg, imageOptions.RotateDegrees)

		width, height := waterImg.Size()
		options.GeoM.Translate(float64(field.X-width/2), float64(field.Y-height/2))
		hexMap.BackgroundSea.DrawImage(waterImg, options)

		r.drawAdjacentPorts(hexMap, board, field, portBgImg)
	})
}

func (r *MapRenderer) drawAdjacentPorts(hexMap *HexMap, board *Board, field *Field, portBgImg []*ebiten.Image) {
	for n := 0; n < 6; n++ {
		neighbor := board.neighborOf(field, n)
		if neighbor == nil || neighbor.Estate != EstatePort {
			continue
		}

		portImg := portBgImg[portImageNum[n]-1]
		portOptions := &ebiten.DrawImageOptions{}
		flipImageMatrix(portOptions, portImg, portFlipH[n], portFlipV[n])

		width, height := portImg.Size()
		portOptions.GeoM.Translate(float64(field.X)-float64(width/2), float64(field.Y)-float64(height/2))
		hexMap.BackgroundSea.DrawImage(portImg, portOptions)
	}
}

func (r *MapRenderer) drawTownBackground(hexMap *HexMap, board *Board) {
	townBgDirtImg := loadImages(fileDir+"/cd_%d.png", 6)
	townBgGrassImg := loadImages(fileDir+"/c_%d.png", 6)

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Estate != EstateTown && field.Estate != EstatePort {
			return
		}
		imageOptions := hexMap.townImageOptions[fieldKey(x, y)]
		r.drawTown(hexMap, field, imageOptions, townBgDirtImg, townBgGrassImg)
	})
}

// drawTown draws the dirt and grass sprites for a town/port tile. Both
// sprites share a single orientation matrix derived from the grass
// image's dimensions.
func (r *MapRenderer) drawTown(
	hexMap *HexMap,
	field *Field,
	imageOptions *LandImageOptions,
	townBgDirtImg []*ebiten.Image,
	townBgGrassImg []*ebiten.Image,
) {
	dirtImg := townBgDirtImg[imageOptions.BgDirtImgIndex]
	grassImg := townBgGrassImg[imageOptions.BgGrassImgIndex]

	options := &ebiten.DrawImageOptions{}
	flipImageMatrix(options, grassImg, imageOptions.FlipH, imageOptions.FlipV)
	rotateImageMatrix(options, grassImg, imageOptions.RotateDegrees)

	width, height := grassImg.Size()
	options.GeoM.Translate(float64(field.X-width/2), float64(field.Y-height/2))
	hexMap.BackgroundDirt.DrawImage(dirtImg, options)
	hexMap.BackgroundGrass.DrawImage(grassImg, options)
}

func (r *MapRenderer) drawTowns(hexMap *HexMap, board *Board) {
	capitalsImg := loadNamedImages([4]string{"capital_red.png", "capital_violet.png", "capital_blue.png", "capital_green.png"})
	cityImage := loadImage(fileDir + "/city.png")
	portImage := loadImage(fileDir + "/port.png")

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		r.drawHex(hexMap.BackgroundTiles, float64(field.X), float64(field.Y))

		switch field.Estate {
		case EstateTown:
			img := cityImage
			if field.Capital >= 0 {
				img = capitalsImg[field.Capital]
			}
			// Sized off cityImage even for capitals, so the icon lines up
			// with the plain city sprite it's replacing.
			r.drawCentered(hexMap.BackgroundTiles, img, cityImage, field)
		case EstatePort:
			r.drawCentered(hexMap.BackgroundTiles, portImage, portImage, field)
		}
	})
}

func (r *MapRenderer) drawCentered(dest *ebiten.Image, img *ebiten.Image, sizingImg *ebiten.Image, field *Field) {
	options := &ebiten.DrawImageOptions{}
	width, height := sizingImg.Size()
	options.GeoM.Translate(float64(field.X)-float64(width/2), float64(field.Y)-float64(height/2))
	dest.DrawImage(img, options)
}

func (r *MapRenderer) drawHex(background *ebiten.Image, xCenter float64, yCenter float64) {
	lineColor := color.RGBA{255, 255, 102, 50}
	drawLine(background, xCenter-12.5, yCenter-20, xCenter-23, yCenter-0, lineColor)
	drawLine(background, xCenter-23, yCenter-0, xCenter-12.5, yCenter+20, lineColor)
	drawLine(background, xCenter-12.5, yCenter+20, xCenter+12.5, yCenter+20, lineColor)
	drawLine(background, xCenter+12.5, yCenter+20, xCenter+23, yCenter+0, lineColor)
	drawLine(background, xCenter+23, yCenter+0, xCenter+12.5, yCenter-20, lineColor)
	drawLine(background, xCenter+12.5, yCenter-20, xCenter-12.5, yCenter-20, lineColor)
}

func drawLine(background *ebiten.Image, x1 float64, y1 float64, x2 float64, y2 float64, lineColor color.RGBA) {
	vector.StrokeLine(background, float32(x1), float32(y1), float32(x2), float32(y2), 0.5, lineColor, true)
}

func (r *MapRenderer) drawTownNames(hexMap *HexMap, board *Board) {
	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Estate != EstateTown && field.Estate != EstatePort {
			return
		}
		textX := field.X - (len(field.TownName) * 3)
		textY := field.Y - int(float64(board.HexHeight)/3)
		// Drawn twice at the same position to fake a bold, more legible font.
		text.Draw(hexMap.BackgroundTiles, field.TownName, hexMap.TextFont, textX, textY, color.White)
		text.Draw(hexMap.BackgroundTiles, field.TownName, hexMap.TextFont, textX, textY, color.White)
	})
}

func (r *MapRenderer) composeBackground(hexMap *HexMap, board *Board) {
	mapView := image.Rect(0, 0,
		int(float64(board.HexWidth)*(float64(board.XMax)/1.35)),
		int(float64(board.HexHeight)*(float64(board.YMax)+0.5)))

	hexMap.Background.DrawImage(hexMap.BackgroundDirt.SubImage(mapView).(*ebiten.Image), nil)
	hexMap.Background.DrawImage(hexMap.BackgroundGrass.SubImage(mapView).(*ebiten.Image), nil)
	hexMap.Background.DrawImage(hexMap.BackgroundSea.SubImage(mapView).(*ebiten.Image), nil)
	hexMap.Background.DrawImage(hexMap.BackgroundTiles, nil)
}

func (r *MapRenderer) drawUI(hexMap *HexMap) {
	redButtonImage := loadImage(fileDir + "/red_button.png")
	randomMapOptions := &ebiten.DrawImageOptions{}
	randomMapOptions.GeoM.Translate(250, 500)
	hexMap.UI.DrawImage(redButtonImage, randomMapOptions)
	text.Draw(hexMap.UI, "Random Map", hexMap.TextFont, 260, 517, color.White)

	grayButtonImage := loadImage(fileDir + "/gray_button.png")
	mapNumberOptions := &ebiten.DrawImageOptions{}
	mapNumberOptions.GeoM.Translate(mapNumberButtonX, mapNumberButtonY)
	hexMap.UI.DrawImage(grayButtonImage, mapNumberOptions)

	text.Draw(hexMap.UI, "Map Number", hexMap.TextFont, 75, 517, color.White)
}

// mapNumberInputHighlightColor outlines the map number button while the
// player is typing a new one, distinguishing it from the plain display
// state without changing the (always black) text color.
var mapNumberInputHighlightColor = color.RGBA{255, 140, 0, 255}

// Position and size of the gray map number button, matching where
// drawUI places it, so the edit-mode outline lines up with it.
const (
	mapNumberButtonX      = 150
	mapNumberButtonY      = 500
	mapNumberButtonWidth  = 98
	mapNumberButtonHeight = 32
)

// DrawMapNumberOverlay draws the map number button's text: the live
// input buffer, with a cursor shown at its actual edit position, while
// the player is typing a new map number, or the current map's number
// otherwise. It's drawn separately from the static UI layer, fresh every
// frame, so keystrokes appear without regenerating the whole map.
func (r *MapRenderer) DrawMapNumberOverlay(dest *ebiten.Image, hexMap *HexMap) {
	label := fmt.Sprintf("%v", hexMap.Board.MapNumber)
	if hexMap.numberInput.Active() {
		digits := hexMap.numberInput.Text()
		cursor := hexMap.numberInput.CursorIndex()
		label = digits[:cursor] + "|" + digits[cursor:]
		vector.StrokeRect(dest,
			mapNumberButtonX, mapNumberButtonY, mapNumberButtonWidth, mapNumberButtonHeight,
			3, mapNumberInputHighlightColor, true)
	}
	text.Draw(dest, label, hexMap.TextFont, 175, 517, color.Black)
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func loadImages(pathFormat string, count int) []*ebiten.Image {
	images := make([]*ebiten.Image, count)
	for i := 0; i < count; i++ {
		images[i] = loadImage(fmt.Sprintf(pathFormat, i+1))
	}
	return images
}

func loadNamedImages(fileNames [4]string) []*ebiten.Image {
	images := make([]*ebiten.Image, len(fileNames))
	for i, name := range fileNames {
		images[i] = loadImage(fileDir + "/" + name)
	}
	return images
}

func flipImageMatrix(options *ebiten.DrawImageOptions, img *ebiten.Image, h int, v int) {
	imgWidth, imgHeight := img.Size()
	width := float64(imgWidth)
	height := float64(imgHeight)

	if h > 0 && v > 0 {
		options.GeoM.Scale(-1, -1)
		options.GeoM.Translate(width, height)
	} else if h > 0 {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(width, 0)
	} else if v > 0 {
		options.GeoM.Scale(1, -1)
		options.GeoM.Translate(0, height)
	}
}

func rotateImageMatrix(options *ebiten.DrawImageOptions, img *ebiten.Image, rotateDegrees int) {
	imgWidth, imgHeight := img.Size()
	width := float64(imgWidth)
	height := float64(imgHeight)

	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Rotate(degreesToRadians(rotateDegrees))
	options.GeoM.Translate(width/2, height/2)
}

func degreesToRadians(degrees int) float64 {
	return (math.Pi / 180) * float64(degrees)
}
