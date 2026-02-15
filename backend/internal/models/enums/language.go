package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Language struct {
	value string
}

var (
	English = Language{"EN"}
	French  = Language{"FR"}
)

func (l Language) String() string {
	return l.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for Language.
func (l *Language) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case English.value:
		*l = English
	case French.value:
		*l = French
	default:
		return fmt.Errorf("invalid language value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Language.
func (l Language) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.value)
}

// Scan implements the sql.Scanner interface for Language.
func (l *Language) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("language cannot be null")
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Language", value)
	}

	*l = Language{value: s}
	return nil
}

// Value implements the driver.Valuer interface for Language.
func (l Language) Value() (driver.Value, error) {
	return l.value, nil
}
