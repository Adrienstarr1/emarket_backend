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

func (h *UserHandler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		UserRequest: request,
		id:          uuid.NewString(),
		created_at:  time.Now(),
		admin:       false,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.UserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.UserRequest.Password = string(hash)

	err = repo.AddUser(r.Context(), h.Pool, user.id, user.UserRequest.Email, user.UserRequest.Age, user.UserRequest.Name, user.UserRequest.Password, user.admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ss, err := auth.CreateSS(user.UserRequest.Name, user.id, user.UserRequest.Email, user.admin)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ss)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	var user UserResponse
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
			user = UserResponse{
				Name:  response.Name,
				Email: response.Email,
				Id:    response.Id,
				Admin: response.Admin,
			}
		}
	}

	if !found {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	ss, err := auth.CreateSS(user.Name, user.Id, user.Email, user.Admin)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
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

	user := UserResponse{
		Id:    responses[0].Id,
		Email: responses[0].Email,
		Name:  responses[0].Name,
		Age:   responses[0].Age,
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

func (h *UserHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Invalid url path", http.StatusBadRequest)
		return
	}

	list, err := repo.ListUsers(r.Context(), h.Pool, "name", name)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
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
			return
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

	if user.Admin {
		acceptedinputs = append(acceptedinputs, admininputs...)
	}

	for key := range instructions {
		if !slices.Contains(acceptedinputs, key) {
			delete(instructions, key)
		}
	}

	err = repo.UpdateUser_V3(r.Context(), h.Pool, user.Id, instructions)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated")

}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.Metadata).(auth.User)

	if err := repo.DeleteUser(r.Context(), h.Pool, user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
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
