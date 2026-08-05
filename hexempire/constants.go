package hexempire

// FieldType classifies terrain: whether a field is land or water.
type FieldType string

const (
	FieldTypeLand  FieldType = "land"
	FieldTypeWater FieldType = "water"
)

// Estate classifies what, if anything, is built on a field.
type Estate string

const (
	EstateNone Estate = ""
	EstateTown Estate = "town"
	EstatePort Estate = "port"
)
