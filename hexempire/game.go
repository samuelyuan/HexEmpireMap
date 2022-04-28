package hexempire

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	backgroundColor = color.NRGBA{0x00, 0x00, 0x00, 0xff}
)

type Game struct {
	hexMap *HexMap

	ScreenWidth  int
	ScreenHeight int
}

func NewGame(textFont font.Face) *Game {
	game := &Game{}
	game.hexMap = NewHexMap(rand.Intn(999999), textFont)
	game.hexMap.generateMap()
	game.ScreenWidth = 800
	game.ScreenHeight = 600
	return game
}

func (game *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if game.hexMap.isMouseCursorOnRandomMapButton(x, y) {
			// Generate new random map
			game.hexMap = NewHexMap(rand.Intn(999999), game.hexMap.TextFont)
			game.hexMap.generateMap()
		}
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	game.hexMap.drawBackground(screen)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}
