package hexempire

import "strconv"

// Field is a single hex cell on the board.
type Field struct {
	FX        int
	FY        int
	X         int
	Y         int
	LandId    int
	Type      FieldType
	Capital   int
	Neighbors [6]*Point2D
	IsLand    bool
	Estate    Estate
	TownName  string
}

// Point2D is a hex-grid coordinate (FX, FY), not a pixel position.
type Point2D struct {
	X int
	Y int
}

func fieldKey(x int, y int) string {
	return "f" + strconv.Itoa(x) + "x" + strconv.Itoa(y)
}
