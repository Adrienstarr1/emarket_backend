package handler

import (
	"e-market/auth"
	"e-market/model"
	"encoding/json"
	"net/http"
)

func (h *Handler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)
	product_id := r.PathValue("product_id")
	if user.Id == "" || product_id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	instructions := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.Add2cart(r.Context(), user.Id, product_id, instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Created")
}

func (h *Handler) DeleteCartHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)
	product_id := r.PathValue("product_id")
	if user.Id == "" || product_id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	instructions := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Service.DeleteCart(r.Context(), user.Id, product_id, instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}
