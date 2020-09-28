package sim

import (
	"testing"
)

func TestCalculateHealed(t *testing.T) {
	const maxHP = 15

	tests := []struct {
		roll    int
		current int
		want    int
	}{
		{roll: 0, current: 10, want: 0},
		{roll: 1, current: 10, want: 1},
		{roll: 5, current: 10, want: 5},
		{roll: 6, current: 10, want: 5},
		{roll: 100, current: 10, want: 5},
		{roll: 5, current: 11, want: 4},
		{roll: 100, current: 1, want: 14},
		{roll: 1, current: maxHP, want: 0},
	}

	for _, test := range tests {
		have := calculateHealed(test.roll, test.current, maxHP)
		if have != test.want {
			t.Errorf("roll=%d current=%d max=%d:\nhave: %d\nwant: %d",
				test.roll, test.current, maxHP, have, test.want)
		}
	}
}
