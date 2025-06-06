package algorithms

import (
	"github.com/IAmRadek/packing/internal/domain/pack"
)

// Allocator defines interface for different algorithms.
type Allocator interface {
	// Allocate returns a map[PackID]PacksUsed to cover demand units, or error.
	Allocate(sizes pack.Sizes, demand int64) (map[pack.ID]pack.Quantity, error)
}
