package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type UnitType struct {
	value string
}

var (
	Volume    = UnitType{"VOLUME"}
	Weight    = UnitType{"WEIGHT"}
	Numeric   = UnitType{"NUMERIC"}
	Piece     = UnitType{"PIECE"}
	Bag       = UnitType{"BAG"}
	Undefined = UnitType{"UNDEFINED"}
)

func (u UnitType) String() string {
	return u.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for UnitType.
func (u *UnitType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case Volume.value:
		*u = Volume
	case Weight.value:
		*u = Weight
	case Numeric.value:
		*u = Numeric
	case Piece.value:
		*u = Piece
	case Bag.value:
		*u = Bag
	case Undefined.value:
		*u = Undefined
	default:
		return fmt.Errorf("invalid unit type value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for UnitType.
func (u UnitType) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.value)
}

// Scan implements the sql.Scanner interface for UnitType.
func (u *UnitType) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("unit type cannot be null")
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into UnitType", value)
	}

	*u = UnitType{value: s}
	return nil
}

// Value implements the driver.Valuer interface for UnitType.
func (u UnitType) Value() (driver.Value, error) {
	return u.value, nil
}
