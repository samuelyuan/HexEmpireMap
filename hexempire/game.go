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
	hexMapId := rand.Intn(999999)
	game.hexMap = NewHexMap(hexMapId, textFont)
	game.hexMap.generateMap()
	game.ScreenWidth = 800
	game.ScreenHeight = 600
	return game
}

func (game *Game) Update() error {
	game.handleClick()
	game.handleMapNumberTyping()
	return nil
}

func (game *Game) handleClick() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	x, y := ebiten.CursorPosition()
	input := &game.hexMap.numberInput

	switch {
	case game.hexMap.isMouseCursorOnRandomMapButton(x, y):
		game.loadMap(rand.Intn(999999))
	case input.Active():
		// Re-clicking the button mid-edit keeps editing as-is; clicking
		// anywhere else cancels it. Either way, don't wipe the buffer.
		if !game.hexMap.isMouseCursorOnMapNumberButton(x, y) {
			input.Cancel()
		}
	case game.hexMap.isMouseCursorOnMapNumberButton(x, y):
		input.Begin(game.hexMap.Board.MapNumber)
	}
}

// handleMapNumberTyping reads keyboard input while the map number button
// is being edited: typed digits insert at the cursor, backspace/delete
// remove around it, arrow/home/end move it, escape cancels, and
// enter/numpad-enter loads the typed map number.
func (game *Game) handleMapNumberTyping() {
	input := &game.hexMap.numberInput
	if !input.Active() {
		return
	}

	for _, r := range ebiten.AppendInputChars(nil) {
		input.InsertDigit(r)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		input.DeleteBefore()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		input.DeleteAfter()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		input.MoveLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		input.MoveRight()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		input.MoveToStart()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		input.MoveToEnd()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		input.Cancel()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) {
		if mapNumber, ok := input.Confirm(); ok {
			game.loadMap(mapNumber)
		}
	}
}

func (game *Game) loadMap(mapNumber int) {
	game.hexMap = NewHexMap(mapNumber, game.hexMap.TextFont)
	game.hexMap.generateMap()
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	game.hexMap.drawBackground(screen)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}
