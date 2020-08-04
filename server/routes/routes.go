package routes

import (
	"auth/controllers"
	"auth/datastore"
	"auth/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Handlers(db datastore.Datastore) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Serve index page on all unhandled routes
	r.HandleFunc("/api/ping", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Print("Pong")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{
			"message": "Pong from the ping route",
		})
	})
	r.HandleFunc("/createUser", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateUser(w, r, db)
	}).Methods("POST")
	r.HandleFunc("/Login", func(writer http.ResponseWriter, request *http.Request) {
		controllers.Login(writer, request, db)
	}).Methods("POST")

	// Authentication that requires the token

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middleware.JwsVerification)
	s.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUsers(w, r, db)
	}).Methods("GET")
	s.HandleFunc("/bar/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUser(w, r, db)
	}).Methods("GET")

	s.HandleFunc("/AddPost", func(writer http.ResponseWriter, request *http.Request) {
		controllers.AddPost(writer, request, db)
	}).Methods("POST")

	s.HandleFunc("/GetPosts", func(writer http.ResponseWriter, request *http.Request) {
		controllers.GetPosts(writer, request, db)
	}).Methods("GET")

	s.HandleFunc("/GetPostsByPriority", func(writer http.ResponseWriter, request *http.Request) {
		controllers.GetPostsByPriority(writer, request, db)
	}).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/build/")))

	return r
}
