package product

import (
	"eboox/auth"
	repo "eboox/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	Id           string    `json:"id"`
	Product_name string    `json:"name"`
	Quantity     int       `json:"quantity"`
	User_Id      string    `json:"user_id"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
	Price        int       `json:"price"`
	Description  string    `json:"description"`
}

type ProductHandler struct {
	Pool *pgxpool.Pool
}

func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	user := r.Context().Value(auth.Metadata).(auth.User)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Error decoding strings", http.StatusInternalServerError)
		return
	}

	product.Id = uuid.New().String()
	product.User_Id = user.Id
	repo.AddProduct(r.Context(), h.Pool, product.Id, product.Product_name, product.Quantity, product.User_Id, product.Price, product.Description)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func (h *ProductHandler) FindProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	response, err := repo.FindProduct_V2(r.Context(), h.Pool, "id", id)
	if err != nil {
		http.Error(w, "ERROR\nCant retrive product", http.StatusInternalServerError)
		return
	}

	product := Product{
		Id:           response[0].Id,
		User_Id:      response[0].User_Id,
		Product_name: response[0].Product_name,
		Quantity:     response[0].Quantity,
		Price:        response[0].Price,
		Description:  response[0].Description,
	}

	json.NewEncoder(w).Encode(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	var instructions map[string]any

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		log.Println(err)
	}

	err = repo.UpdateProduct_V2(r.Context(), h.Pool, id, instructions)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Updated")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	if err := repo.DeleteProduct(r.Context(), h.Pool, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Deleted")
	w.WriteHeader(http.StatusOK)
}
