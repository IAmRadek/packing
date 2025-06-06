package handlers

import (
	"net/http"

	"github.com/IAmRadek/packing/internal/templates"
)

type IndexHandler struct {
	render *templates.Templates
}

func NewIndexHandler(render *templates.Templates) *IndexHandler {
	return &IndexHandler{
		render: render,
	}
}

func (h *IndexHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	h.render.Render(w, r, "index", nil)
}
