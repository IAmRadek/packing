package allocation

import (
	"context"
	"fmt"

	"github.com/IAmRadek/packing/internal/algorithms"
	"github.com/IAmRadek/packing/internal/domain/pack"
)

type Repo interface {
	GetInventory(ctx context.Context, sku string) (*pack.Inventory, error)
}

type Service struct {
	repo      Repo
	allocator algorithms.Allocator
}

func NewService(repo Repo, algo algorithms.Allocator) *Service {
	return &Service{
		repo:      repo,
		allocator: algo,
	}
}

func (s *Service) Compute(ctx context.Context, sku string, quantity int64) (pack.Allocations, error) {
	inv, err := s.repo.GetInventory(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("getting inventory: %w", err)
	}

	sizes := inv.AvailableSizes()

	dist, err := s.allocator.Allocate(sizes, quantity)
	if err != nil {
		return nil, fmt.Errorf("allocating: %w", err)
	}

	out := make(pack.Allocations, 0, len(dist))
	for id, qty := range dist {
		size, _ := sizes.ByID(id)
		out = append(out, pack.Allocation{
			Size:     size,
			Quantity: qty,
		})
	}

	return out, nil
}
