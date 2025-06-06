package dp

import (
	"fmt"
	"testing"
)

func TestAllocate(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []int64
		quantity int64
		exp      map[int64]int64
	}{
		{
			name:     "edge_case",
			sizes:    []int64{23, 31, 53},
			quantity: 500_000,
			exp: map[int64]int64{
				23: 2,
				31: 7,
				53: 9429,
			},
		},
		{
			name:     "simple",
			sizes:    []int64{250, 500, 1000, 2000, 5000},
			quantity: 1,
			exp: map[int64]int64{
				250: 1,
			},
		},
		{
			name:     "tricky",
			sizes:    []int64{250, 500, 1000, 2000, 5000},
			quantity: 251,
			exp: map[int64]int64{
				500: 1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			allocation := Allocate(tc.sizes, tc.quantity)

			if err := cmp(allocation, tc.exp); err != nil {
				t.Errorf("Allocate() mismatch:\n%s", err.Error())
			}
		})
	}
}

func cmp(a, b map[int64]int64) error {
	if len(a) != len(b) {
		return fmt.Errorf("len(a) != len(b)")
	}
	for k, v := range a {
		if b[k] != v {
			return fmt.Errorf("got[%d](%d) != want[%d](%d)", k, v, k, b[k])
		}
	}

	return nil
}

func BenchmarkAllocate(b *testing.B) {
	sizes := []int64{250, 500, 1000, 2000, 5000}

	demands := []int64{
		1, 500, 2500, 10000, 50000, 100000, 500000, 1000000, 5000000, 10000000,
	}

	for _, d := range demands {
		b.Run(fmt.Sprintf("demand:%d", d), func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				Allocate(sizes, d)
			}

			b.Logf("%d*%d demand calculated in %s", d, b.N, b.Elapsed())
		})
	}
}
