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
	A     any
	Child *child
}

var (
	swS                          = []string{"S"}
	swI                          = []string{"I"}
	swU                          = []string{"U"}
	swF                          = []string{"F"}
	swA                          = []string{"A"}
	swChild                      = []string{"Child"}
	swChildGrandChild            = []string{"Child", "GrandChild"}
	swChildGrandChildI           = []string{"Child", "GrandChild", "I"}
	swChildGrandChildPossiblyNil = []string{"Child", "GrandChild", "PossiblyNil"}

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
