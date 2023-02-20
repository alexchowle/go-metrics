package metrics

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func TestGenericGauge(t *testing.T) {
	type test[T constraints.Ordered] struct {
		setTwice bool
		input    T
		want     T
	}

	int64Tests := map[string]test[int64]{
		"simple int64 set":    {input: 1, want: 1},
		"setting int64 twice": {setTwice: true, input: 1, want: 1},
	}
	for name, tc := range int64Tests {
		t.Run(name, func(t *testing.T) {
			g := NewGenericGauge[int64]()
			g.Update(tc.input)
			if tc.setTwice {
				g.Update(tc.input)
			}
			got := g.Value()
			if got != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}

	float64Tests := map[string]test[float64]{
		"simple float64 set":    {input: 1.0, want: 1.0},
		"setting float64 twice": {setTwice: true, input: 1.0, want: 1.0},
	}
	for name, tc := range float64Tests {
		t.Run(name, func(t *testing.T) {
			g := NewGenericGauge[float64]()
			g.Update(tc.input)
			if tc.setTwice {
				g.Update(tc.input)
			}
			got := g.Value()
			if got != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}

	stringTests := map[string]test[string]{
		"simple string set":    {input: "hello", want: "hello"},
		"setting string twice": {setTwice: true, input: "hello", want: "hello"},
	}
	for name, tc := range stringTests {
		t.Run(name, func(t *testing.T) {
			g := NewGenericGauge[string]()
			g.Update(tc.input)
			if tc.setTwice {
				g.Update(tc.input)
			}
			got := g.Value()
			if got != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
