package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestAvatar_String(t *testing.T) {
	tests := []struct {
		name     string
		avatar   Avatar
		expected string
	}{
		{"Avatar1", Avatar1, "/static/assets/avatar1.jpg"},
		{"Avatar2", Avatar2, "/static/assets/avatar2.jpg"},
		{"Avatar3", Avatar3, "/static/assets/avatar3.jpg"},
		{"Empty avatar", Avatar{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.avatar.String(); actual != tt.expected {
				t.Errorf("Avatar.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestAvatar_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  Avatar
		expectErr bool
	}{
		{"Valid Avatar1", `"/static/assets/avatar1.jpg"`, Avatar1, false},
		{"Valid Avatar2", `"/static/assets/avatar2.jpg"`, Avatar2, false},
		{"Valid Avatar3", `"/static/assets/avatar3.jpg"`, Avatar3, false},
		{"Null value", `null`, Avatar{}, false},
		{"Empty string", `""`, Avatar{}, false},
		{"Invalid value", `"invalid.jpg"`, Avatar{}, true},
		{"Invalid JSON", `invalid`, Avatar{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var avatar Avatar
			err := json.Unmarshal([]byte(tt.jsonData), &avatar)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && avatar != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", avatar, tt.expected)
			}
		})
	}
}

func TestAvatar_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		avatar   Avatar
		expected string
	}{
		{"Avatar1", Avatar1, `"/static/assets/avatar1.jpg"`},
		{"Avatar2", Avatar2, `"/static/assets/avatar2.jpg"`},
		{"Avatar3", Avatar3, `"/static/assets/avatar3.jpg"`},
		{"Empty avatar", Avatar{}, `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.avatar)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, expected %v", string(data), tt.expected)
			}
		})
	}
}

func TestAvatar_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  Avatar
		expectErr bool
	}{
		{"Valid string", "/static/assets/avatar1.jpg", Avatar{value: "/static/assets/avatar1.jpg"}, false},
		{"Another valid string", "/static/assets/avatar2.jpg", Avatar{value: "/static/assets/avatar2.jpg"}, false},
		{"Nil value", nil, Avatar{}, false},
		{"Invalid type", 123, Avatar{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var avatar Avatar
			err := avatar.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && avatar.value != tt.expected.value {
				t.Errorf("Scan() = %v, expected %v", avatar.value, tt.expected.value)
			}
		})
	}
}

func TestAvatar_Value(t *testing.T) {
	tests := []struct {
		name     string
		avatar   Avatar
		expected driver.Value
	}{
		{"Avatar1", Avatar1, "/static/assets/avatar1.jpg"},
		{"Avatar2", Avatar2, "/static/assets/avatar2.jpg"},
		{"Avatar3", Avatar3, "/static/assets/avatar3.jpg"},
		{"Empty avatar returns nil", Avatar{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.avatar.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if val != tt.expected {
				t.Errorf("Value() = %v, expected %v", val, tt.expected)
			}
		})
	}
}
