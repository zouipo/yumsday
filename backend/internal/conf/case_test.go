package conf

import (
	"slices"
	"testing"
)

func TestSplitCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected []string
	}{
		{
			name:     "Valid",
			value:    "oneTwo",
			expected: []string{"one", "Two"},
		},
		{
			name:     "Single word",
			value:    "word",
			expected: []string{"word"},
		},
		{
			name:     "Acronym",
			value:    "anACROnym",
			expected: []string{"an", "ACR", "Onym"},
		},
		{
			name:     "With number",
			value:    "one2three",
			expected: []string{"one", "2", "three"},
		},
		{
			name:     "With number and upper case",
			value:    "one2THREEFour",
			expected: []string{"one", "2", "THREE", "Four"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := splitCamelCase(test.value)
			if slices.Compare(test.expected, got) != 0 {
				t.Errorf("Expected %v, got %v", test.expected, got)
			}
		})
	}
}

func TestToConstantCase(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{
			value:    "oneTwo",
			expected: "ONE_TWO",
		},
		{
			value:    "oneTWOThree",
			expected: "ONE_TWO_THREE",
		},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got := toConstantCase(test.value)
			if test.expected != got {
				t.Errorf("Expected %v, got %v", test.expected, got)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{
			value:    "oneTwo",
			expected: "one-two",
		},
		{
			value:    "oneTWOThree",
			expected: "one-two-three",
		},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got := toKebabCase(test.value)
			if test.expected != got {
				t.Errorf("Expected %v, got %v", test.expected, got)
			}
		})
	}
}
