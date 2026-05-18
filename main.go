package main

import (
	"context"
	"e-market/auth"
	"e-market/handler"
	"e-market/repo"
	"e-market/service"
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
	pool, err := repo.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	Handler := handler.Handler{
		Service: service.Service{
			Repo: repo.Repo{
				Pool: pool,
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/signin", Handler.SignupHandler)
	mux.HandleFunc("POST /user/login", Handler.LoginHandler)
	mux.HandleFunc("POST /product/create", auth.Auth(Handler.CreateProductHandler))
	mux.HandleFunc("PATCH /product/update/{id}", auth.Auth(Handler.UpdateProductHandler))
	mux.HandleFunc("GET /product/{id}", auth.Auth(Handler.GetProduct))
	mux.HandleFunc("GET /user/{id}/{info}", auth.Auth(Handler.UserInfoHandler))
	mux.HandleFunc("GET /user/{id}", auth.Auth(Handler.UserInfoHandler))
	mux.HandleFunc("GET /user/list/{name}", auth.Auth(Handler.ListUsersHandler))
	mux.HandleFunc("PATCH /user/update", auth.Auth(Handler.UpdateUserhandler))
	mux.HandleFunc("DELETE /user/delete", auth.Auth(Handler.DeleteUserHandler))
	mux.HandleFunc("POST /user/addcart/{product_id}", auth.Auth(Handler.AddToCartHandler))
	mux.HandleFunc("PATCH /user/makeadmin/{id}", auth.Auth(Handler.MakeUserAdminHandler))

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
