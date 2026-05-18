package handler

import (
	"e-market/auth"
	"e-market/model"
	"e-market/repo"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := model.User{
		Name:       request.Name,
		Email:      request.Email,
		Password:   request.Password,
		Age:        request.Age,
		Id:         uuid.NewString(),
		Created_at: time.Now(),
		Admin:      false,
	}

	ss, err := h.Service.Signup(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ss)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error decoding strings", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Name:       request.Name,
		Email:      request.Email,
		Password:   request.Password,
		Age:        request.Age,
		Id:         uuid.NewString(),
		Created_at: time.Now(),
		Admin:      false,
	}

	ss, err := h.Service.Login(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ss)
}

func (h *Handler) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	user, err := h.Service.UserInfo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info := r.PathValue("info")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	switch info {
	case "name":
		json.NewEncoder(w).Encode(user.Name)
	case "email":
		json.NewEncoder(w).Encode(user.Email)
	case "age":
		json.NewEncoder(w).Encode(user.Age)
	default:
		json.NewEncoder(w).Encode(user)
	}

}

func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var repo *repo.Repo
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	list, err := repo.ListUsers(r.Context(), "name", name)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) UpdateUserhandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)
	var instructions map[string]any

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		log.Println(err)
	}

	err = h.Service.UpdateUser(r.Context(), user, instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated")

}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)

	if err := h.Service.Repo.DeleteUser(r.Context(), user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

func (h *Handler) MakeUserAdminHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)
	id := r.PathValue("id")

	err := h.Service.Makeadmin(r.Context(), user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListUsersProductsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(model.User)

	response, err := h.Service.Repo.FindProduct_V2(r.Context(), "user_id", user.Id)
	if err != nil {
		http.Error(w, "ERROR\nCant retrive product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
