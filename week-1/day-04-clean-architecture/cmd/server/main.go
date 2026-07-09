// Package main is the entry point for the user API server.
// It wires together all layers (repository, service, handler) and starts the HTTP server.
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/meedaycodes/day04-clean-architecture/internal/handler"
	"github.com/meedaycodes/day04-clean-architecture/internal/repository"
	"github.com/meedaycodes/day04-clean-architecture/internal/service"
)

func main() {

	repo := repository.NewInMemoryUserRepository()
	serv := service.NewUserService(repo)
	servHandler := handler.NewUserHandler(serv)

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	combinedHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			servHandler.GetAllUsers(w, req)
		} else if req.Method == http.MethodPost {
			servHandler.CreateUser(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/users/", servHandler.GetUserByID)
	http.HandleFunc("/users", combinedHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
