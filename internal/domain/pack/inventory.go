package pack

type Inventory struct {
	sku   string
	packs Sizes

	// TODO: keep track of Stock
}

func (i *Inventory) SKU() string {
	return i.sku
}

func NewInventory(sku string, sizes Sizes) *Inventory {
	return &Inventory{
		sku:   sku, // TODO: slugify
		packs: sizes,
	}
}

func (i *Inventory) Update(sizes Sizes) {
	i.packs = sizes
}

func (i *Inventory) AvailableSizes() Sizes {
	return i.packs
}
