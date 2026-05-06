package main

import (
	"context"
	"eboox/auth"
	"eboox/order"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := repository.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	UserHandler := user.UserHandler{
		Pool: pool,
	}
	ProductHandler := product.ProductHandler{
		Pool: pool,
	}
	OrderHandler := order.OrderHandler{
		Pool: pool,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/signin", UserHandler.SiginHandler)
	mux.HandleFunc("POST /user/login", UserHandler.LoginHandler)
	mux.HandleFunc("POST /product/create", auth.Auth(ProductHandler.CreateProductHandler))
	mux.HandleFunc("PATCH /product/update/{id}", auth.Auth(ProductHandler.UpdateProductHandler))
	mux.HandleFunc("GET /product/{id}", auth.Auth(ProductHandler.FindProduct))
	mux.HandleFunc("GET /user/{id}/{info}", auth.Auth(UserHandler.UserInfoHandler))
	mux.HandleFunc("GET /user/{id}", auth.Auth(UserHandler.UserInfoHandler))
	mux.HandleFunc("GET /user/list/{name}", auth.Auth(UserHandler.ListUsersHandler))
	mux.HandleFunc("PATCH /user/update", auth.Auth(UserHandler.UpdateUserhandler))
	mux.HandleFunc("DELETE /user/delete", auth.Auth(UserHandler.DeleteUserHandler))
	mux.HandleFunc("POST /user/addcart/{product_id}", auth.Auth(OrderHandler.AddToCartHandler))
	mux.HandleFunc("PATCH /user/makeadmin/{id}", auth.Auth(UserHandler.MakeUserAdminHandler))

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
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()
	<-stopchan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer pool.Close()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server shutting Failed %v", err)
	}
	log.Println("Server shutdown sucessfull")
}
