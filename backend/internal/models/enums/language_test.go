package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestLanguage_String(t *testing.T) {
	tests := []struct {
		name     string
		language Language
		expected string
	}{
		{"English", English, "EN"},
		{"French", French, "FR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.language.String(); actual != tt.expected {
				t.Errorf("Language.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestLanguage_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  Language
		expectErr bool
	}{
		{"Valid English", `"EN"`, English, false},
		{"Valid French", `"FR"`, French, false},
		{"Invalid value", `"ES"`, Language{}, true},
		{"Invalid JSON", `invalid`, Language{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var language Language
			err := json.Unmarshal([]byte(tt.jsonData), &language)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && language != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", language, tt.expected)
			}
		})
	}
}

func TestLanguage_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		language Language
		expected string
	}{
		{"English", English, `"EN"`},
		{"French", French, `"FR"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.language)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, expected %v", string(data), tt.expected)
			}
		})
	}
}

func TestLanguage_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  Language
		expectErr bool
	}{
		{"Valid string", "EN", Language{value: "EN"}, false},
		{"Another valid string", "FR", Language{value: "FR"}, false},
		{"Nil value", nil, Language{}, true},
		{"Invalid type", 123, Language{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var language Language
			err := language.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && language.value != tt.expected.value {
				t.Errorf("Scan() = %v, expected %v", language.value, tt.expected.value)
			}
		})
	}
}

func TestLanguage_Value(t *testing.T) {
	tests := []struct {
		name     string
		language Language
		expected driver.Value
	}{
		{"English", English, "EN"},
		{"French", French, "FR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.language.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if val != tt.expected {
				t.Errorf("Value() = %v, expected %v", val, tt.expected)
			}
		})
	}
}
