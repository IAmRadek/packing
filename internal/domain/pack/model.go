package pack

import (
	"fmt"
)

type ID string

type Size struct {
	ID       ID
	Capacity int64
	Label    string
}

type Sizes []Size

func NewSizes(capacities []int64, labels []string) (Sizes, error) {
	if len(capacities) != len(labels) {
		return nil, fmt.Errorf("capacities and labels must have the same length")
	}
	if len(capacities) == 0 {
		return nil, fmt.Errorf("capacities and labels must have at least one element")
	}

	if hasDuplicates(capacities) {
		return nil, fmt.Errorf("capacities must not have duplicates")
	}

	out := make(Sizes, len(capacities))
	for i, c := range capacities {
		out[i] = Size{
			ID:       ID(labels[i]),
			Capacity: c,
			Label:    labels[i],
		}
	}

	return out, nil
}

func (s Sizes) Combine(other Sizes) (Sizes, error) {
	if len(s) == 0 {
		return other, nil
	}
	if len(other) == 0 {
		return s, nil
	}

	out := make(Sizes, 0, len(s)+len(other))
	for _, s := range s {
		out = append(out, s)
	}
	for _, s := range other {
		out = append(out, s)

	}

	if hasDuplicates(out.Capacities()) {
		return nil, fmt.Errorf("capacities must not have duplicates")
	}

	return out, nil
}

func hasDuplicates(capacities []int64) bool {
	seen := make(map[int64]struct{})
	for _, c := range capacities {
		if _, ok := seen[c]; ok {
			return true
		}
		seen[c] = struct{}{}
	}
	return false
}

func (s Sizes) ByID(id ID) (Size, bool) {
	for _, s := range s {
		if s.ID == id {
			return s, true
		}
	}
	return Size{}, false
}

func (s Sizes) ByCapacity(cap int64) (Size, bool) {
	for _, s := range s {
		if s.Capacity == cap {
			return s, true
		}
	}
	return Size{}, false
}

func (s Sizes) Capacities() []int64 {
	out := make([]int64, 0, len(s))
	for _, s := range s {
		out = append(out, s.Capacity)
	}
	return out
}

type Quantity int64

type Allocation struct {
	Size     Size
	Quantity Quantity
}

type Allocations []Allocation

func (a Allocations) SumItems() int64 {
	out := int64(0)
	for _, a := range a {
		out += int64(a.Quantity) * a.Size.Capacity
	}
	return out
}

func (a Allocations) SumPacks() int64 {
	out := int64(0)
	for _, a := range a {
		out += int64(a.Quantity)
	}
	return out
}
