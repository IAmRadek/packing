package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/IAmRadek/packing/internal/app/allocation"
	"github.com/IAmRadek/packing/internal/app/inventory"
	"github.com/IAmRadek/packing/internal/domain/pack"
	"github.com/IAmRadek/packing/internal/templates"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type InventoryHandler struct {
	invSrv   *inventory.Service
	allocSrv *allocation.Service
	render   *templates.Templates
	dec      *schema.Decoder
}

func NewInventoryHandler(
	invSrv *inventory.Service,
	allocSrv *allocation.Service,
	render *templates.Templates,
	dec *schema.Decoder,
) *InventoryHandler {
	return &InventoryHandler{
		invSrv:   invSrv,
		allocSrv: allocSrv,
		render:   render,
		dec:      dec,
	}
}

func (h *InventoryHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	invs, err := h.invSrv.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Render(w, r, "inventory_list", map[string]interface{}{
		"inventories": invs,
	})
}

type InventoryCreateRequest struct {
	Name       string   `schema:"name"`
	Labels     []string `schema:"pack_name[]"`
	Quantities []int64  `schema:"pack_quantity[]"`
}

type InventoryCreateResponse struct {
	Error string
	Name  string
}

func (h *InventoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	resp := InventoryCreateResponse{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req InventoryCreateRequest

		if err := h.dec.Decode(&req, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			resp.Error = "name is required"
			h.render.Render(w, r, "inventory_create", resp)
			return
		}

		if strings.EqualFold(req.Name, "create") {
			resp.Error = "name cannot be 'create'"
			h.render.Render(w, r, "inventory_create", resp)
			return
		}

		sizes, err := pack.NewSizes(req.Quantities, req.Labels)
		if err != nil {
			resp.Error = err.Error()
			h.render.Render(w, r, "inventory_create", resp)
			return
		}

		if err := h.invSrv.Create(r.Context(), sanitize(req.Name), sizes); err != nil {
			resp.Error = err.Error()
			h.render.Render(w, r, "inventory_create", resp)
			return
		}
	}

	h.render.Render(w, r, "inventory_create", resp)
}

type InventoryGetRequest struct {
	Demand int64 `schema:"demand"`
}

type InventoryGetResponse struct {
	Inventory   *pack.Inventory
	Demand      int64
	Allocations pack.Allocations
}

func (h *InventoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["sku"] == "" {
		http.Error(w, "sku is required", http.StatusBadRequest)
		return
	}

	resp := InventoryGetResponse{}

	inv, err := h.invSrv.Get(r.Context(), vars["sku"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Inventory = inv

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req InventoryGetRequest

		if err := h.dec.Decode(&req, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		packs, err := h.allocSrv.Compute(r.Context(), inv.SKU(), req.Demand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.Demand = req.Demand
		resp.Allocations = packs
	}

	h.render.Render(w, r, "inventory_get", resp)
}

type InventoryUpdateRequest struct {
	SKU           string   `schema:"sku"`
	Labels        []string `schema:"label[]"`
	Capacities    []int64  `schema:"capacity[]"`
	NewLabels     []string `schema:"new_label[]"`
	NewCapacities []int64  `schema:"new_capacity[]"`
}

func (h *InventoryHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["sku"] == "" {
		http.Error(w, "sku is required", http.StatusBadRequest)
		return
	}

	inv, err := h.invSrv.Get(r.Context(), vars["sku"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req InventoryUpdateRequest

		if err := h.dec.Decode(&req, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sizes, err := pack.NewSizes(req.Capacities, req.Labels)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var newSizes pack.Sizes
		if len(req.NewLabels) > 0 {
			var err error
			newSizes, err = pack.NewSizes(req.NewCapacities, req.NewLabels)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		allSizes, err := sizes.Combine(newSizes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.invSrv.Update(r.Context(), vars["sku"], allSizes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.Redirect(w, r, "/inventory/"+inv.SKU(), http.StatusFound)
}

func (h *InventoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["sku"] == "" {
		http.Error(w, "sku is required", http.StatusBadRequest)
		return
	}

	inv, err := h.invSrv.Get(r.Context(), vars["sku"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.invSrv.Delete(r.Context(), inv.SKU()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/inventory", http.StatusFound)
}

func sanitize(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	sanitized := reg.ReplaceAllString(input, "")

	reg = regexp.MustCompile(`\s+`)
	sanitized = reg.ReplaceAllString(sanitized, " ")

	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}
