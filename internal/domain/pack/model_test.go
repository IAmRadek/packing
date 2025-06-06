package pack

import (
	"reflect"
	"testing"
)

func TestNewSizes(t *testing.T) {
	tests := []struct {
		name       string
		capacities []int64
		labels     []string
		want       Sizes
		wantErr    bool
	}{
		{
			name:       "valid input",
			capacities: []int64{10, 20, 30},
			labels:     []string{"small", "medium", "large"},
			want: Sizes{
				{ID: "small", Capacity: 10, Label: "small"},
				{ID: "medium", Capacity: 20, Label: "medium"},
				{ID: "large", Capacity: 30, Label: "large"},
			},
			wantErr: false,
		},
		{
			name:       "different lengths",
			capacities: []int64{10, 20},
			labels:     []string{"small"},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty input",
			capacities: []int64{},
			labels:     []string{},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "duplicate capacities",
			capacities: []int64{10, 10, 30},
			labels:     []string{"small", "medium", "large"},
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSizes(tt.capacities, tt.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSizes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSizes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizes_Combine(t *testing.T) {
	s1 := Sizes{{ID: "small", Capacity: 10, Label: "small"}}
	s2 := Sizes{{ID: "large", Capacity: 20, Label: "large"}}
	sDuplicate := Sizes{{ID: "medium", Capacity: 10, Label: "medium"}}

	tests := []struct {
		name    string
		s       Sizes
		other   Sizes
		want    Sizes
		wantErr bool
	}{
		{
			name:  "combine non-empty sizes",
			s:     s1,
			other: s2,
			want: Sizes{
				{ID: "small", Capacity: 10, Label: "small"},
				{ID: "large", Capacity: 20, Label: "large"},
			},
			wantErr: false,
		},
		{
			name:    "first empty",
			s:       Sizes{},
			other:   s2,
			want:    s2,
			wantErr: false,
		},
		{
			name:    "second empty",
			s:       s1,
			other:   Sizes{},
			want:    s1,
			wantErr: false,
		},
		{
			name:    "duplicate capacities",
			s:       s1,
			other:   sDuplicate,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Combine(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Combine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizes_ByID(t *testing.T) {
	sizes := Sizes{
		{ID: "small", Capacity: 10, Label: "small"},
		{ID: "medium", Capacity: 20, Label: "medium"},
	}

	tests := []struct {
		name      string
		id        ID
		want      Size
		wantFound bool
	}{
		{
			name:      "existing id",
			id:        "small",
			want:      Size{ID: "small", Capacity: 10, Label: "small"},
			wantFound: true,
		},
		{
			name:      "non-existing id",
			id:        "large",
			want:      Size{},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := sizes.ByID(tt.id)
			if found != tt.wantFound {
				t.Errorf("ByID() found = %v, wantFound %v", found, tt.wantFound)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizes_ByCapacity(t *testing.T) {
	sizes := Sizes{
		{ID: "small", Capacity: 10, Label: "small"},
		{ID: "medium", Capacity: 20, Label: "medium"},
	}

	tests := []struct {
		name      string
		capacity  int64
		want      Size
		wantFound bool
	}{
		{
			name:      "existing capacity",
			capacity:  10,
			want:      Size{ID: "small", Capacity: 10, Label: "small"},
			wantFound: true,
		},
		{
			name:      "non-existing capacity",
			capacity:  30,
			want:      Size{},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := sizes.ByCapacity(tt.capacity)
			if found != tt.wantFound {
				t.Errorf("ByCapacity() found = %v, wantFound %v", found, tt.wantFound)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllocations_SumItems(t *testing.T) {
	tests := []struct {
		name        string
		allocations Allocations
		want        int64
	}{
		{
			name: "multiple allocations",
			allocations: Allocations{
				{Size: Size{Capacity: 10}, Quantity: 2},
				{Size: Size{Capacity: 20}, Quantity: 3},
			},
			want: 80, // (10 * 2) + (20 * 3)
		},
		{
			name:        "empty allocations",
			allocations: Allocations{},
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.allocations.SumItems(); got != tt.want {
				t.Errorf("SumItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllocations_SumPacks(t *testing.T) {
	tests := []struct {
		name        string
		allocations Allocations
		want        int64
	}{
		{
			name: "multiple allocations",
			allocations: Allocations{
				{Quantity: 2},
				{Quantity: 3},
			},
			want: 5, // 2 + 3
		},
		{
			name:        "empty allocations",
			allocations: Allocations{},
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.allocations.SumPacks(); got != tt.want {
				t.Errorf("SumPacks() = %v, want %v", got, tt.want)
			}
		})
	}
}
