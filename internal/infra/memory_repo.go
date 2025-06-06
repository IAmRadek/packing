package infra

import (
	"context"
	"fmt"
	"sync"

	"github.com/IAmRadek/packing/internal/domain/pack"
)

type MemoryRepo struct {
	rw *sync.RWMutex
	m  map[string]*pack.Inventory
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		rw: &sync.RWMutex{},
		m:  make(map[string]*pack.Inventory),
	}
}

func (m *MemoryRepo) ListInventories(ctx context.Context) ([]*pack.Inventory, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	out := make([]*pack.Inventory, 0, len(m.m))
	for _, inv := range m.m {
		out = append(out, inv)
	}
	return out, nil
}

func (m *MemoryRepo) DeleteInventory(ctx context.Context, sku string) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	delete(m.m, sku)
	return nil
}

func (m *MemoryRepo) Save(ctx context.Context, inv *pack.Inventory) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	if _, ok := m.m[inv.SKU()]; ok {
		return fmt.Errorf("inventory already exists for sku: %s", inv.SKU())
	}

	m.m[inv.SKU()] = inv
	return nil
}

func (m *MemoryRepo) GetInventory(ctx context.Context, sku string) (*pack.Inventory, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	inv, ok := m.m[sku]
	if !ok {
		return nil, fmt.Errorf("inventory not found for sku: %s", sku)
	}
	return inv, nil
}
