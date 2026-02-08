package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// AppTheme represents the theme preference of the application. It can be one of the predefined themes.
type AppTheme struct {
	value string
}

var (
	Light  = AppTheme{"LIGHT"}
	Dark   = AppTheme{"DARK"}
	System = AppTheme{"SYSTEM"}
)

func (a AppTheme) String() string {
	return a.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for AppTheme.
func (a *AppTheme) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case Light.value:
		*a = Light
	case Dark.value:
		*a = Dark
	case System.value:
		*a = System
	default:
		return fmt.Errorf("invalid app theme value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for AppTheme.
func (a AppTheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.value)
}

// Scan implements the sql.Scanner interface for AppTheme.
func (a *AppTheme) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("app theme cannot be null")
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into AppTheme", value)
	}

	*a = AppTheme{value: s}
	return nil
}

// Value implements the driver.Valuer interface for AppTheme.
func (a AppTheme) Value() (driver.Value, error) {
	return a.value, nil
}
