package enums

type UnitType string

const (
	Volume    UnitType = "VOLUME"
	Weight    UnitType = "WEIGHT"
	Numeric   UnitType = "NUMERIC"
	Piece     UnitType = "PIECE"
	Bag       UnitType = "BAG"
	Undefined UnitType = "UNDEFINED"
)
