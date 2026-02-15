package enums

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Avatar represents a user's avatar image. It can be one of the predefined avatars or null.
type Avatar struct {
	value string
}

var (
	Avatar1 = Avatar{"/static/assets/avatar1.jpg"}
	Avatar2 = Avatar{"/static/assets/avatar2.jpg"}
	Avatar3 = Avatar{"/static/assets/avatar3.jpg"}
)

func (a Avatar) String() string {
	return a.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for Avatar.
// It handles null values and empty strings by setting the Avatar to an empty value.
func (a *Avatar) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*a = Avatar{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Handle empty string as null
	if s == "" {
		*a = Avatar{}
		return nil
	}

	switch s {
	case Avatar1.value:
		*a = Avatar1
	case Avatar2.value:
		*a = Avatar2
	case Avatar3.value:
		*a = Avatar3
	default:
		return fmt.Errorf("invalid avatar value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Avatar.
// It returns null for empty Avatar values.
func (a Avatar) MarshalJSON() ([]byte, error) {
	if a.value == "" {
		return []byte("null"), nil
	}
	return json.Marshal(a.value)
}

// Scan implements the sql.Scanner interface for Avatar.
func (a *Avatar) Scan(value interface{}) error {
	if value == nil {
		*a = Avatar{}
		return nil
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Avatar", value)
	}

	*a = Avatar{value: s}
	return nil
}

// Value implements the driver.Valuer interface for Avatar.
func (a Avatar) Value() (driver.Value, error) {
	if a.value == "" {
		return nil, nil
	}
	return a.value, nil
}
