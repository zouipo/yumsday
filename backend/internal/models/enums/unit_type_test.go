package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestUnitType_String(t *testing.T) {
	tests := []struct {
		name     string
		unitType UnitType
		expected string
	}{
		{"Volume", Volume, "VOLUME"},
		{"Weight", Weight, "WEIGHT"},
		{"Numeric", Numeric, "NUMERIC"},
		{"Piece", Piece, "PIECE"},
		{"Bag", Bag, "BAG"},
		{"Undefined", Undefined, "UNDEFINED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.unitType.String(); actual != tt.expected {
				t.Errorf("UnitType.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestUnitType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  UnitType
		expectErr bool
	}{
		{"Valid Volume", `"VOLUME"`, Volume, false},
		{"Valid Weight", `"WEIGHT"`, Weight, false},
		{"Valid Numeric", `"NUMERIC"`, Numeric, false},
		{"Valid Piece", `"PIECE"`, Piece, false},
		{"Valid Bag", `"BAG"`, Bag, false},
		{"Valid Undefined", `"UNDEFINED"`, Undefined, false},
		{"Invalid value", `"INVALID"`, UnitType{}, true},
		{"Invalid JSON", `invalid`, UnitType{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unitType UnitType
			err := json.Unmarshal([]byte(tt.jsonData), &unitType)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && unitType != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", unitType, tt.expected)
			}
		})
	}
}

func TestUnitType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		unitType UnitType
		expected string
	}{
		{"Volume", Volume, `"VOLUME"`},
		{"Weight", Weight, `"WEIGHT"`},
		{"Numeric", Numeric, `"NUMERIC"`},
		{"Piece", Piece, `"PIECE"`},
		{"Bag", Bag, `"BAG"`},
		{"Undefined", Undefined, `"UNDEFINED"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.unitType)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, expected %v", string(data), tt.expected)
			}
		})
	}
}

func TestUnitType_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  UnitType
		expectErr bool
	}{
		{"Valid string", "VOLUME", UnitType{value: "VOLUME"}, false},
		{"Another valid string", "WEIGHT", UnitType{value: "WEIGHT"}, false},
		{"Nil value", nil, UnitType{}, true},
		{"Invalid type", 123, UnitType{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unitType UnitType
			err := unitType.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && unitType.value != tt.expected.value {
				t.Errorf("Scan() = %v, expected %v", unitType.value, tt.expected.value)
			}
		})
	}
}

func TestUnitType_Value(t *testing.T) {
	tests := []struct {
		name     string
		unitType UnitType
		expected driver.Value
	}{
		{"Volume", Volume, "VOLUME"},
		{"Weight", Weight, "WEIGHT"},
		{"Numeric", Numeric, "NUMERIC"},
		{"Piece", Piece, "PIECE"},
		{"Bag", Bag, "BAG"},
		{"Undefined", Undefined, "UNDEFINED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.unitType.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if val != tt.expected {
				t.Errorf("Value() = %v, expected %v", val, tt.expected)
			}
		})
	}
}
