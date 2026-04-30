package main

import (
	"context"
	product "eboox/products"
	"eboox/repository"
	user "eboox/users"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load()
	connString  = os.Getenv("DB_URL")
	pool, _     = repository.Connect(context.Background())
	UserHandler = user.UserHandler{
		Pool: pool,
	}
	ProductHandler = product.ProductHandler{
		Pool: pool,
	}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/signin", UserHandler.SiginHandler)
	mux.HandleFunc("POST /user/login", UserHandler.LoginHandler)
	mux.HandleFunc("POST /product/create", user.Auth(ProductHandler.CreateProductHandler))
	mux.HandleFunc("PUT /product/update/{id}", user.Auth(ProductHandler.UpdateProductHandler))
	mux.HandleFunc("GET /product/{id}", user.Auth(ProductHandler.FindProduct))
	mux.HandleFunc("GET /user/{id}/{info}", user.Auth(UserHandler.UserInfoHandler))
	mux.HandleFunc("GET /user/{id}", user.Auth(UserHandler.UserInfoHandler))
	mux.HandleFunc("GET /user/list/{name}", user.Auth(UserHandler.ListUsersHandler))
	mux.HandleFunc("PUT /user/update/{id}", user.Auth(UserHandler.UpdateUserhandler))
	mux.HandleFunc("DELETE /user/delete/{id}", user.Auth(UserHandler.DeleteUserHandler))
	mux.HandleFunc("POST /user/addcart/{user_id}/{product_id}", user.Auth(UserHandler.AddToCartHandler))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, os.Interrupt)
	go func() {
		log.Println("Listening on ...8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()
	<-stopchan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer pool.Close()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server shutting Failed %v", err)
	}
	log.Println("Server shutdown sucessfull")
}
