package utils

import "testing"

type grandChild struct {
	I           int
	PossiblyNil any
}

type child struct {
	GrandChild *grandChild
}

type parent struct {
	S     string
	I     int
	U     uint
	F     float32
	Child *child
}

type otherKind struct {
	I uint
}

var (
	swS                = []string{"S"}
	swI                = []string{"I"}
	swU                = []string{"U"}
	swF                = []string{"F"}
	swChild            = []string{"Child"}
	swChildGrandChildI = []string{"Child", "GrandChild", "I"}

	a = parent{S: "a"}
	b = parent{S: "b"}

	minusOne = parent{I: -1}
	minusTwo = parent{I: -2}

	one = parent{U: 1}
	two = parent{U: 2}

	half      = parent{F: 0.5}
	minusHalf = parent{F: -0.5}

	nestedOne = parent{
		Child: &child{
			GrandChild: &grandChild{
				I: 1,
			},
		},
	}
	nestedTwo = parent{
		Child: &child{
			GrandChild: &grandChild{
				I: 2,
			},
		},
	}

	nilPointer = parent{Child: nil}

	oneOtherKind = otherKind{I: 1}
)

func TestCompareFieldsByName(t *testing.T) {
	tests := []struct {
		name      string
		smallest  parent
		greatest  parent
		sortWords []string
	}{
		{
			name:      "compare_string",
			smallest:  a,
			greatest:  b,
			sortWords: swS,
		},
		{
			name:      "compare_int",
			smallest:  minusTwo,
			greatest:  minusOne,
			sortWords: swI,
		},
		{
			name:      "compare_uint",
			smallest:  one,
			greatest:  two,
			sortWords: swU,
		},
		{
			name:      "compare_float",
			smallest:  minusHalf,
			greatest:  half,
			sortWords: swF,
		},
		{
			name:      "compare_struct",
			smallest:  nestedOne,
			greatest:  nestedTwo,
			sortWords: swChildGrandChildI,
		},
		{
			name:      "compare_nil",
			smallest:  nilPointer,
			greatest:  nestedTwo,
			sortWords: swChildGrandChildI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_less_than", func(t *testing.T) {
			actual := compareFieldsByName(tt.smallest, tt.greatest, tt.sortWords, false)
			if !actual {
				t.Error("expected comparison to return true")
			}
		})

		t.Run(tt.name+"_more_than", func(t *testing.T) {
			actual := compareFieldsByName(tt.greatest, tt.smallest, tt.sortWords, false)
			if actual {
				t.Error("expected comparison to return false")
			}
		})

		t.Run(tt.name+"_descending", func(t *testing.T) {
			actual := compareFieldsByName(tt.smallest, tt.greatest, tt.sortWords, true)
			if actual {
				t.Error("expected comparison to return false")
			}
		})
	}
}

func TestCompareFieldsByName_TwoNils(t *testing.T) {
	actualAscending := compareFieldsByName(nilPointer, nilPointer, swChildGrandChildI, false)
	actualDescending := compareFieldsByName(nilPointer, nilPointer, swChildGrandChildI, true)

	if !actualAscending {
		t.Error("comparison of two nil pointers in ascending order should always return true")
	}

	if actualDescending {
		t.Error("comparison of two nil pointers in descending order should always return false")
	}
}

func TestCompareFieldsByName_DifferentKinds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	// need to convert to any here because the compiler
	// won't let us pass two values of different types
	compareFieldsByName(any(one), any(oneOtherKind), swI, false)
}

func TestCompareFieldsByName_CannotCompareStructs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	compareFieldsByName(nestedOne, nestedTwo, swChild, false)
}
