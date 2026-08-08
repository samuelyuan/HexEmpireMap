package hexempire

import (
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
)

// HexMap composes map generation and rendering: it holds the rendered
// image layers and the currently generated Board, and delegates the two
// concerns to MapGenerator and MapRenderer respectively.
type HexMap struct {
	Background      *ebiten.Image
	BackgroundDirt  *ebiten.Image
	BackgroundGrass *ebiten.Image
	BackgroundSea   *ebiten.Image
	BackgroundTiles *ebiten.Image
	UI              *ebiten.Image
	Board           *Board
	TextFont        font.Face

	generator              *MapGenerator
	backgroundImageOptions map[string]*LandImageOptions
	waterImageOptions      map[string]*WaterImageOptions
	townImageOptions       map[string]*LandImageOptions

	numberInput MapNumberInput
}

func NewHexMap(mapNumber int, textFont font.Face) *HexMap {
	hexMap := &HexMap{}
	hexMap.generator = NewMapGenerator(mapNumber)
	hexMap.Background = ebiten.NewImage(800, 600)
	hexMap.BackgroundDirt = ebiten.NewImage(800, 600)
	hexMap.BackgroundGrass = ebiten.NewImage(800, 600)
	hexMap.BackgroundSea = ebiten.NewImage(800, 600)
	hexMap.BackgroundTiles = ebiten.NewImage(800, 600)
	hexMap.UI = ebiten.NewImage(800, 600)
	hexMap.TextFont = textFont
	return hexMap
}

func (hexMap *HexMap) generateBoard() {
	generated := hexMap.generator.Generate()
	hexMap.Board = generated.Board
	hexMap.backgroundImageOptions = generated.BackgroundImageOptions
	hexMap.waterImageOptions = generated.WaterImageOptions
	hexMap.townImageOptions = generated.TownImageOptions
}

func (hexMap *HexMap) generateMap() {
	hexMap.generateBoard()
	NewMapRenderer().Render(hexMap)
}

func (hexMap *HexMap) isMouseCursorOnRandomMapButton(x int, y int) bool {
	return x >= 255 && x <= 345 && y >= 500 && y <= 530
}

// isMouseCursorOnMapNumberButton reports whether (x, y) is over the map
// number button, which can be clicked to type in a specific map number.
func (hexMap *HexMap) isMouseCursorOnMapNumberButton(x int, y int) bool {
	return x >= 155 && x <= 240 && y >= 500 && y <= 530
}

func (hexMap *HexMap) drawBackground(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(25, 25)
	screen.DrawImage(hexMap.Background, options)
	screen.DrawImage(hexMap.UI, nil)

	// Drawn fresh every frame (instead of baked into the static UI image)
	// so typed digits show up immediately without regenerating the map.
	NewMapRenderer().DrawMapNumberOverlay(screen, hexMap)
}
