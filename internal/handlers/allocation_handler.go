package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IAmRadek/packing/internal/app/allocation"
)

type AllocationHandler struct {
	srv *allocation.Service
}

func NewAllocationHandler(srv *allocation.Service) *AllocationHandler {
	return &AllocationHandler{
		srv: srv,
	}
}

type AllocateRequest struct {
	Sku      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

func (h *AllocationHandler) HandleAllocate(w http.ResponseWriter, r *http.Request) {
	var req AllocateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Sku == "" {
		http.Error(w, "sku is required", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be positive", http.StatusBadRequest)
	}

	packs, err := h.srv.Compute(r.Context(), req.Sku, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(packs)
	return
}
