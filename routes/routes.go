package routes

import (
	"auth/controllers"
	"auth/datastore"
	"auth/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func Handlers(db datastore.Datastore) * mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/createUser", func(w http.ResponseWriter, r* http.Request) {
		controllers.CreateUser(w, r, db)
	}).Methods("POST")
	r.HandleFunc("/Login", func(writer http.ResponseWriter, request *http.Request) {
		controllers.Login(writer, request, db)
	}).Methods("POST")

	// Authentication that requires the token

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middleware.JwsVerification)
	s.HandleFunc("/foo", func (w http.ResponseWriter, r* http.Request){
		controllers.GetUsers(w,r,db)
	}).Methods("GET")
	s.HandleFunc("/bar/{id}", func (w http.ResponseWriter, r* http.Request){
		controllers.GetUser(w,r,db)
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

	return r
}
