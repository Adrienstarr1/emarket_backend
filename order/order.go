package order

import (
	"eboox/auth"
	repo "eboox/repository"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	Pool *pgxpool.Pool
}

func (h *OrderHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)
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

	instructions["user_id"] = user.Id
	instructions["product_id"] = product_id

	res, err := repo.FindCart(r.Context(), h.Pool, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exist := func() bool {
		for _, r := range res {
			if product_id == r.Product_id {
				return true
			}
		}
		return false
	}()

	switch exist {
	case true:
		err = repo.UpdateCart(r.Context(), h.Pool, user.Id, product_id, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case false:
		err = repo.Add2cart(r.Context(), h.Pool, instructions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Created")
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteCartHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)
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

	res, err := repo.FindCart(r.Context(), h.Pool, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exist, available := func() (bool, bool) {
		for _, r := range res {
			if product_id == r.Product_id && r.Quantity > 0 {
				return true, true
			} else if product_id == r.Product_id && r.Quantity <= 0 {
				return true, false
			}
		}
		return false, false
	}()

	switch exist {
	case true:
		switch available {
		case true:
			err = repo.UpdateCart(r.Context(), h.Pool, user.Id, product_id, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case false:
			err = repo.DeleteCart(r.Context(), h.Pool, user.Id, product_id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case false:
		http.Error(w, "cart not found", http.StatusBadRequest)
		return
	}
}
