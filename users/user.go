package user

import (
	"eboox/auth"
	repo "eboox/repository"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) SiginHandler(w http.ResponseWriter, r *http.Request) {
	var instructions map[string]any
	err := json.NewDecoder(r.Body).Decode(&instructions)

	acceptedinputs := []string{
		"name",
		"age",
		"email",
		"password",
	}

	for key := range instructions {
		if !slices.Contains(acceptedinputs, key) {
			delete(instructions, key)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	instructions["id"] = uuid.New().String()
	instructions["admin"] = false
	instructions["created_at"] = time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(instructions["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instructions["password"] = string(hash)

	err = repo.AddUser_V2(r.Context(), h.Pool, instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ss := auth.CreateSS(instructions["name"].(string), instructions["id"].(string), instructions["email"].(string), instructions["admin"].(bool))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ss)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	var user User
	var found bool
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding strings", http.StatusInternalServerError)
		return
	}
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	responses, err := repo.ListUsers(r.Context(), h.Pool, "email", input.Email)
	if err != nil {
		http.Error(w, "Problem retriving user:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(responses) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for _, response := range responses {
		err = bcrypt.CompareHashAndPassword([]byte(response.Password), []byte(input.Password))
		if err == nil {
			found = true
			user = User{
				Name:  response.Name,
				Email: response.Email,
				id:    response.Id,
				admin: response.Admin,
			}
		}
	}

	if !found {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	ss := auth.CreateSS(user.Name, user.id, user.Email, user.admin)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ss)
}

func (h *UserHandler) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	responses, err := repo.ListUsers(r.Context(), h.Pool, "id", id)
	if err != nil {
		http.Error(w, "Problem retriving user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(responses) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user := User{
		id:    responses[0].Id,
		Email: responses[0].Email,
		Name:  responses[0].Name,
		Age:   responses[0].Age,
	}

	info := r.PathValue("info")
	w.Header().Set("Content-Type", "application/json")
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

	w.WriteHeader(http.StatusOK)

}

func (h *UserHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if "name" == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	list, err := repo.ListUsers(r.Context(), h.Pool, "name", name)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdateUserhandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)
	var instructions map[string]any

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		log.Println(err)
	}

	password, exist := func() (string, bool) {
		for key, value := range instructions {
			if key == "password" {
				password := value.(string)
				return password, true
			}
		}
		return "", false
	}()

	if exist {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		instructions["password"] = string(hash)
	}

	acceptedinputs := []string{
		"name",
		"age",
		"email",
		"password",
	}
	admininputs := []string{
		"id",
		"created_at",
	}

	for key := range instructions {
		if !slices.Contains(acceptedinputs, key) && !user.Admin {
			delete(instructions, key)
		} else if !slices.Contains(admininputs, key) && user.Admin {
			delete(instructions, key)
		}
	}

	err = repo.UpdateUser_V3(r.Context(), h.Pool, user.Id, instructions)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Updated")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)

	if err := repo.DeleteUser(r.Context(), h.Pool, user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Deleted")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) MakeUserAdminHandler(w http.ResponseWriter, r *http.Request) {
	instructions := make(map[string]any)
	user := r.Context().Value(auth.Metadata).(auth.User)
	if !user.Admin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	instructions["admin"] = true

	if err := repo.UpdateUser_V3(r.Context(), h.Pool, id, instructions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ListUsersProductsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)

	response, err := repo.FindProduct_V2(r.Context(), h.Pool, "user_id", user.Id)
	if err != nil {
		http.Error(w, "ERROR\nCant retrive product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
