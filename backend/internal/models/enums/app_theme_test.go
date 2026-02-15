package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestAppTheme_String(t *testing.T) {
	tests := []struct {
		name     string
		theme    AppTheme
		expected string
	}{
		{"Light theme", Light, "LIGHT"},
		{"Dark theme", Dark, "DARK"},
		{"System theme", System, "SYSTEM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.theme.String(); actual != tt.expected {
				t.Errorf("AppTheme.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestAppTheme_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  AppTheme
		expectErr bool
	}{
		{"Valid Light", `"LIGHT"`, Light, false},
		{"Valid Dark", `"DARK"`, Dark, false},
		{"Valid System", `"SYSTEM"`, System, false},
		{"Invalid value", `"INVALID"`, AppTheme{}, true},
		{"Invalid JSON", `invalid`, AppTheme{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var theme AppTheme
			err := json.Unmarshal([]byte(tt.jsonData), &theme)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && theme != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", theme, tt.expected)
			}
		})
	}
}

func TestAppTheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		theme    AppTheme
		expected string
	}{
		{"Light theme", Light, `"LIGHT"`},
		{"Dark theme", Dark, `"DARK"`},
		{"System theme", System, `"SYSTEM"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.theme)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, expected %v", string(data), tt.expected)
			}
		})
	}
}

func TestAppTheme_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  AppTheme
		expectErr bool
	}{
		{"Valid string", "LIGHT", AppTheme{value: "LIGHT"}, false},
		{"Another valid string", "DARK", AppTheme{value: "DARK"}, false},
		{"Nil value", nil, AppTheme{}, true},
		{"Invalid type", 123, AppTheme{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var theme AppTheme
			err := theme.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && theme.value != tt.expected.value {
				t.Errorf("Scan() = %v, expected %v", theme.value, tt.expected.value)
			}
		})
	}
}

func TestAppTheme_Value(t *testing.T) {
	tests := []struct {
		name     string
		theme    AppTheme
		expected driver.Value
	}{
		{"Light theme", Light, "LIGHT"},
		{"Dark theme", Dark, "DARK"},
		{"System theme", System, "SYSTEM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.theme.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if val != tt.expected {
				t.Errorf("Value() = %v, expected %v", val, tt.expected)
			}
		})
	}
}
