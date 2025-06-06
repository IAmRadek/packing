package dp

import (
	"math"

	"github.com/IAmRadek/packing/internal/domain/pack"
)

type Allocator struct{}

func (a Allocator) Allocate(sizes pack.Sizes, demand int64) (map[pack.ID]pack.Quantity, error) {
	dist := Allocate(sizes.Capacities(), demand)

	out := make(map[pack.ID]pack.Quantity, len(dist))
	for k, v := range dist {
		s, _ := sizes.ByCapacity(k)
		out[s.ID] = pack.Quantity(v)
	}

	return out, nil
}

// Allocate tries to distribute `demand` into packs of `sizes`
func Allocate(sizes []int64, demand int64) map[int64]int64 {
	// Note: find the greatest common divisor so we can shrink the search space.
	g := gcd(sizes)
	for i := range sizes {
		sizes[i] /= g
	}
	demand = (demand + g - 1) / g // ceil so total ≥ original demand/g

	maxS := max(sizes)
	limit := demand + maxS

	// Note: the following section calculates reachable targets given the sizes.
	reachable := make([]bool, limit+1)
	reachable[0] = true
	for _, s := range sizes {
		for t := s; t <= limit; t++ {
			reachable[t] = reachable[t] || reachable[t-s]
		}
	}

	// Note: searching for the first possible reachable target.
	target := int64(-1)
	for t := demand; t <= limit; t++ {
		if reachable[t] {
			target = t
			break
		}
	}
	// Note: if not found, no solution is possible.
	if target == -1 {
		return nil
	}

	// Note: the following section is calculating possible packs and optimized the run.
	// Adds up packs so we don't have to do it again.
	const inf = math.MaxInt32
	packs := make([]int64, target+1)
	prev := make([]int64, target+1)
	run := make([]int64, target+1)
	for i := range packs {
		packs[i] = inf
	}
	packs[0] = 0
	for _, s := range sizes {
		for t := s; t <= target; t++ {
			if packs[t-s]+1 < packs[t] {
				packs[t] = packs[t-s] + 1
				prev[t] = s
				if prev[t-s] == s {
					run[t] = run[t-s] + 1
				} else {
					run[t] = 1
				}
			}
		}
	}

	// Note: going back from the target and collecting the output.
	out := make(map[int64]int64)
	for t := target; t > 0; {
		s := prev[t]
		k := run[t]
		out[s*g] += k // restore the original unit size from gcd
		t -= s * k
	}

	return out
}

func gcd(xs []int64) int64 {
	calc := func(a, b int64) int64 {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	g := xs[0]
	for _, x := range xs[1:] {
		g = calc(g, x)
	}
	return g
}

func max(xs []int64) int64 {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}
