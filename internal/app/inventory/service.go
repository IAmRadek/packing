package inventory

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/IAmRadek/packing/internal/domain/pack"
)

type Repo interface {
	ListInventories(ctx context.Context) ([]*pack.Inventory, error)
	GetInventory(ctx context.Context, sku string) (*pack.Inventory, error)
	DeleteInventory(ctx context.Context, sku string) error
	Save(ctx context.Context, inv *pack.Inventory) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) List(ctx context.Context) ([]*pack.Inventory, error) {
	lst, err := s.repo.ListInventories(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing inventories: %w", err)
	}

	slices.SortFunc(lst, func(a, b *pack.Inventory) int {
		return strings.Compare(a.SKU(), b.SKU())
	})

	return lst, nil
}

func (s *Service) Create(ctx context.Context, sku string, sizes []pack.Size) error {
	inv := pack.NewInventory(sku, sizes)
	return s.repo.Save(ctx, inv)
}

func (s *Service) Get(ctx context.Context, sku string) (*pack.Inventory, error) {
	return s.repo.GetInventory(ctx, sku)
}

func (s *Service) Update(ctx context.Context, sku string, sizes []pack.Size) error {
	inv, err := s.repo.GetInventory(ctx, sku)
	if err != nil {
		return fmt.Errorf("getting inventory: %w", err)
	}

	inv.Update(sizes)

	return s.repo.Save(ctx, inv)
}

func (s *Service) Delete(ctx context.Context, sku string) error {
	return s.repo.DeleteInventory(ctx, sku)
}
