package user

import (
	repo "eboox/repository"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	MySigningKey = []byte(os.Getenv("MY_SECRET_KEY"))
)

type User struct {
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Id         string    `json:"id"`
	Created_at time.Time `json:"created_at"`
}

type UserHandler struct {
	Pool *pgxpool.Pool
}

type MyClaims struct {
	User User
	jwt.RegisteredClaims
}

func createSS(user User) string {
	claims := MyClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(28 * time.Hour)),
			Issuer:    "Me!!!!!!!!!!!!",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(MySigningKey)

	return ss
}

func (h *UserHandler) SiginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Error decoding Data", http.StatusInternalServerError)
		return
	}

	user.Id = uuid.New().String()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	if err = repo.AddUser(r.Context(), h.Pool, user.Name, user.Age, user.Email, user.Password, user.Id); err != nil {
		http.Error(w, "Error storig user", http.StatusInternalServerError)
	}

	ss := createSS(user)

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
				Name:       response.Name,
				Age:        response.Age,
				Email:      response.Email,
				Id:         response.Id,
				Created_at: response.Created_at,
			}
		}
	}

	if !found {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	ss := createSS(user)
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
		Id:    responses[0].Id,
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

	for key, value := range instructions {
		if key == "password" {
			value, err = bcrypt.GenerateFromPassword([]byte(value.(string)), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	err = repo.UpdateUser_V3(r.Context(), h.Pool, id, instructions)
	if err != nil {
		http.Error(w, "ERROR:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Updated")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	if err := repo.DeleteUser(r.Context(), h.Pool, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Deleted")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	product_id := r.PathValue("product_id")
	if user_id == "" || product_id == "" {
		http.Error(w, "invalid/empty id in path", http.StatusBadRequest)
		return
	}

	instructions := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	instructions["user_id"] = user_id
	instructions["product_id"] = product_id

	res, err := repo.FindCart(r.Context(), h.Pool, user_id)
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
		err = repo.UpdateCart(r.Context(), h.Pool, user_id, product_id)
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
