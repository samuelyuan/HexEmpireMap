package hexempire

import (
	"math"
)

type Pathfinder struct {
}

func NewPathfinder() *Pathfinder {
	return &Pathfinder{}
}

type Tile struct {
	Field     *Field
	Parent    *Tile
	DistCost  float64
	TotalCost float64
}

func (pathfinder *Pathfinder) findPath(board *Board, startField *Field, endField *Field, avoidEstate []string, avoidWater bool) []*Field {
	if startField == nil || endField == nil {
		return nil
	}

	if startField.Type == "water" {
		avoidWater = false
	}

	tiles := make([]*Tile, 0)
	path := make([]*Tile, 0)
	tileLength := make(map[string]int)
	pathFieldLength := make(map[string]int)

	tiles = append(tiles, &Tile{Field: startField})
	tiles[len(tiles)-1].TotalCost = 0
	moveCost := []float64{5, 5, 5, 5, 5, 5}
	for (len(path) == 0 || (len(path) > 0 && path[len(path)-1].Field != endField)) && len(tiles) > 0 {
		currentTile := tiles[0]
		tiles = tiles[1:]
		for neighborNum := 0; neighborNum < 6; neighborNum++ {
			neighbor := board.getNeighborField(currentTile.Field, neighborNum)
			if pathfinder.canWalk(currentTile.Field, neighbor, avoidEstate, avoidWater) ||
				neighbor == endField {
				newTile := &Tile{}
				newTile.Field = neighbor
				distance := pathfinder.getDistance(newTile.Field, endField)
				newTile.Parent = currentTile
				newTile.DistCost = moveCost[neighborNum] + distance
				newTile.TotalCost = currentTile.TotalCost + moveCost[neighborNum]

				key := pathfinder.getFieldStrKey(newTile.Field)
				_, isInPathFieldLength := pathFieldLength[key]
				if !isInPathFieldLength {
					_, isInTileLength := tileLength[key]
					if !isInTileLength {
						tileLength[key] = len(tiles)
						tiles = append(tiles, newTile)
					}
				} else if path[pathFieldLength[key]].TotalCost > newTile.TotalCost {
					path[pathFieldLength[key]] = newTile
				}
			}
		}
		pathFieldLength[pathfinder.getFieldStrKey(currentTile.Field)] = len(path)
		path = append(path, currentTile)
		if len(tiles) > 0 {
			tileForSwap := 0
			for tileNum := 1; tileNum < len(tiles); tileNum++ {
				if tiles[tileNum].DistCost < tiles[tileForSwap].DistCost {
					tileForSwap = tileNum
				}
			}
			temp := tiles[0]
			tiles[0] = tiles[tileForSwap]
			tiles[tileForSwap] = temp
		}
	}
	if len(tiles) == 0 {
		return nil
	}
	finalPath := make([]*Field, 0)
	pathIndex := len(path) - 1
	for len(finalPath) == 0 || (len(finalPath) > 0 && finalPath[len(finalPath)-1] != startField) {
		finalPath = append(finalPath, path[pathIndex].Field)
		// Added to prevent undefined error
		if path[pathIndex].Parent == nil {
			break
		}
		pathIndex = pathFieldLength[pathfinder.getFieldStrKey(path[pathIndex].Parent.Field)]
	}
	pathfinder.reverseArray(finalPath)
	return finalPath
}

func (pathfinder *Pathfinder) canWalk(a *Field, b *Field, avoidEstate []string, avoidWater bool) bool {
	if a == nil || b == nil {
		return false
	}
	for n := 0; n < len(avoidEstate); n++ {
		if b.Estate == avoidEstate[n] {
			return false
		}
	}
	if !avoidWater {
		return true
	}
	if a.Type == "water" && b.Type == "water" {
		return true
	}
	if a.Type == "land" && b.Type == "land" {
		return true
	}
	if a.Type == "water" && b.Type == "land" {
		return true
	}
	if b.Type == "water" && a.Estate == "port" {
		return true
	}
	return false
}

func (pathFinder *Pathfinder) getFieldStrKey(field *Field) string {
	return getFieldKey(field.FX, field.FY)
}

func (pathFinder *Pathfinder) getDistance(a *Field, b *Field) float64 {
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

func (pathfinder *Pathfinder) reverseArray(arr []*Field) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
