package hexempire

// Board holds the full grid of fields plus the derived groupings
// (land masses, towns, capitals) produced by map generation.
type Board struct {
	MapNumber       int
	XMax            int
	YMax            int
	HexWidth        int
	HexHeight       int
	Fields          map[string]*Field
	LandCount       int
	LandGroups      [][]*Field
	Towns           []*Field
	PartiesCapitals []*Field
	TownNames       []string
}

func NewBoard() *Board {
	return &Board{
		XMax:            20,
		YMax:            11,
		HexWidth:        50,
		HexHeight:       40,
		Fields:          make(map[string]*Field),
		LandCount:       0,
		LandGroups:      make([][]*Field, 0),
		Towns:           make([]*Field, 0),
		PartiesCapitals: make([]*Field, 4),
	}
}

func (board *Board) fieldAt(x int, y int) *Field {
	return board.Fields[fieldKey(x, y)]
}

func (board *Board) neighborOf(field *Field, direction int) *Field {
	neighborLocation := field.Neighbors[direction]
	if neighborLocation == nil {
		return nil
	}
	return board.fieldAt(neighborLocation.X, neighborLocation.Y)
}

func (board *Board) locationIfExists(x int, y int) *Point2D {
	if _, ok := board.Fields[fieldKey(x, y)]; ok {
		return &Point2D{X: x, Y: y}
	}
	return nil
}

// isCapitalCorner reports whether (x, y) is one of the four fixed
// corners reserved for each party's starting capital.
func (board *Board) isCapitalCorner(x int, y int) bool {
	return (x == 1 && y == 1) ||
		(x == board.XMax-2 && y == 1) ||
		(x == board.XMax-2 && y == board.YMax-2) ||
		(x == 1 && y == board.YMax-2)
}

// findNeighbors computes and stores the six adjacent grid coordinates for a
// field, accounting for the horizontal offset of odd-numbered columns.
func (board *Board) findNeighbors(field *Field) {
	field.Neighbors = [6]*Point2D{}
	fx := field.FX
	fy := field.FY
	if fx%2 == 0 {
		field.Neighbors[0] = board.locationIfExists(fx+1, fy)
		field.Neighbors[1] = board.locationIfExists(fx, fy+1)
		field.Neighbors[2] = board.locationIfExists(fx-1, fy)
		field.Neighbors[3] = board.locationIfExists(fx-1, fy-1)
		field.Neighbors[4] = board.locationIfExists(fx, fy-1)
		field.Neighbors[5] = board.locationIfExists(fx+1, fy-1)
	} else {
		field.Neighbors[0] = board.locationIfExists(fx+1, fy+1)
		field.Neighbors[1] = board.locationIfExists(fx, fy+1)
		field.Neighbors[2] = board.locationIfExists(fx-1, fy+1)
		field.Neighbors[3] = board.locationIfExists(fx-1, fy)
		field.Neighbors[4] = board.locationIfExists(fx, fy-1)
		field.Neighbors[5] = board.locationIfExists(fx+1, fy)
	}
}

// forEachField visits every (x, y) coordinate in the board's grid in
// row-major order, matching the iteration order generation relies on.
func (board *Board) forEachField(visit func(x int, y int)) {
	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			visit(x, y)
		}
	}
}
