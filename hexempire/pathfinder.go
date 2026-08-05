package hexempire

import "math"

// pathfinderMoveCost is the cost of moving into any adjacent hex; the
// grid has no terrain-dependent movement cost, so it's a flat constant.
const pathfinderMoveCost = 5.0

type Pathfinder struct{}

func NewPathfinder() *Pathfinder {
	return &Pathfinder{}
}

// Tile is a node visited while searching for a path: the field it sits
// on, the tile it was reached from, and the costs used to rank it.
type Tile struct {
	Field     *Field
	Parent    *Tile
	DistCost  float64
	TotalCost float64
}

// findPath searches from startField to endField, preferring low-cost
// tiles (DistCost = move cost + heuristic distance to the goal). Fields
// whose Estate is in avoidEstate are impassable, as is water when
// avoidWater is set (unless the path departs from a water start, or
// crosses from a port).
func (pathfinder *Pathfinder) findPath(board *Board, startField *Field, endField *Field, avoidEstate []Estate, avoidWater bool) []*Field {
	if startField == nil || endField == nil {
		return nil
	}
	if startField.Type == FieldTypeWater {
		avoidWater = false
	}

	search := newPathSearch(board, startField, endField, avoidEstate, avoidWater)
	for len(search.open) > 0 && !search.reachedEnd() {
		search.expandNextTile()
	}

	if len(search.open) == 0 {
		return nil
	}
	return search.reconstructPath(startField)
}

// pathSearch holds the mutable state of a single findPath run: an open
// list of frontier tiles awaiting expansion, and a closed list of tiles
// already expanded, each indexed by field key for O(1) lookup.
type pathSearch struct {
	board       *Board
	endField    *Field
	avoidEstate []Estate
	avoidWater  bool

	open      []*Tile
	openIndex map[string]int

	closed      []*Tile
	closedIndex map[string]int
}

func newPathSearch(board *Board, startField *Field, endField *Field, avoidEstate []Estate, avoidWater bool) *pathSearch {
	return &pathSearch{
		board:       board,
		endField:    endField,
		avoidEstate: avoidEstate,
		avoidWater:  avoidWater,
		open:        []*Tile{{Field: startField}},
		openIndex:   make(map[string]int),
		closedIndex: make(map[string]int),
	}
}

func (search *pathSearch) reachedEnd() bool {
	return len(search.closed) > 0 && search.closed[len(search.closed)-1].Field == search.endField
}

func (search *pathSearch) expandNextTile() {
	currentTile := search.open[0]
	search.open = search.open[1:]

	for direction := 0; direction < 6; direction++ {
		search.considerNeighbor(currentTile, direction)
	}

	search.closedIndex[fieldKey(currentTile.Field.FX, currentTile.Field.FY)] = len(search.closed)
	search.closed = append(search.closed, currentTile)

	search.moveClosestOpenTileToFront()
}

func (search *pathSearch) considerNeighbor(currentTile *Tile, direction int) {
	neighbor := search.board.neighborOf(currentTile.Field, direction)
	if !canWalk(currentTile.Field, neighbor, search.avoidEstate, search.avoidWater) && neighbor != search.endField {
		return
	}

	newTile := &Tile{
		Field:     neighbor,
		Parent:    currentTile,
		DistCost:  pathfinderMoveCost + getDistance(neighbor, search.endField),
		TotalCost: currentTile.TotalCost + pathfinderMoveCost,
	}

	key := fieldKey(neighbor.FX, neighbor.FY)
	if closedPos, inClosed := search.closedIndex[key]; inClosed {
		if search.closed[closedPos].TotalCost > newTile.TotalCost {
			search.closed[closedPos] = newTile
		}
		return
	}
	if _, inOpen := search.openIndex[key]; inOpen {
		return
	}
	search.openIndex[key] = len(search.open)
	search.open = append(search.open, newTile)
}

// moveClosestOpenTileToFront swaps the open tile with the lowest DistCost
// into position 0, so the next expandNextTile call picks it up.
func (search *pathSearch) moveClosestOpenTileToFront() {
	if len(search.open) == 0 {
		return
	}
	closest := 0
	for i := 1; i < len(search.open); i++ {
		if search.open[i].DistCost < search.open[closest].DistCost {
			closest = i
		}
	}
	search.open[0], search.open[closest] = search.open[closest], search.open[0]
}

// reconstructPath walks the closed list backward from the goal via each
// tile's Parent link, then reverses the result into start-to-goal order.
func (search *pathSearch) reconstructPath(startField *Field) []*Field {
	finalPath := make([]*Field, 0)
	tile := search.closed[len(search.closed)-1]
	for {
		finalPath = append(finalPath, tile.Field)
		if tile.Field == startField || tile.Parent == nil {
			break
		}
		parentKey := fieldKey(tile.Parent.Field.FX, tile.Parent.Field.FY)
		tile = search.closed[search.closedIndex[parentKey]]
	}
	reverseFields(finalPath)
	return finalPath
}

func canWalk(a *Field, b *Field, avoidEstate []Estate, avoidWater bool) bool {
	if a == nil || b == nil {
		return false
	}
	for _, estate := range avoidEstate {
		if b.Estate == estate {
			return false
		}
	}
	if !avoidWater {
		return true
	}
	if a.Type == FieldTypeWater && b.Type == FieldTypeWater {
		return true
	}
	if a.Type == FieldTypeLand && b.Type == FieldTypeLand {
		return true
	}
	if a.Type == FieldTypeWater && b.Type == FieldTypeLand {
		return true
	}
	if b.Type == FieldTypeWater && a.Estate == EstatePort {
		return true
	}
	return false
}

func getDistance(a *Field, b *Field) float64 {
	acx := a.FX * 5
	bcx := b.FX * 5
	var acy, bcy int
	if a.FX%2 == 0 {
		acy = a.FY * 10
	} else {
		acy = (a.FY * 10) + 5
	}
	if b.FX%2 == 0 {
		bcy = b.FY * 10
	} else {
		bcy = (b.FY * 10) + 5
	}
	return math.Sqrt(math.Pow(float64(acx-bcx), 2) + math.Pow(float64(acy-bcy), 2))
}

func reverseFields(arr []*Field) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
