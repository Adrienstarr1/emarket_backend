package handler

import (
	"e-market/auth"
	"e-market/model"
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	user := r.Context().Value(auth.Metadata).(model.User)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Error decoding strings", http.StatusInternalServerError)
		return
	}

	err = h.Service.CreateProduct(r.Context(), product, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	product, err := h.Service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}

func (h *Handler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}
	var instructions map[string]any
	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Service.UpdateProduct(r.Context(), id, instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated")

}

func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	if err := h.Service.Repo.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}
