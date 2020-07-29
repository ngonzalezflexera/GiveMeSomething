package controllers

import (
	"auth/datastore"
	"auth/middleware"
	"auth/model"
	"auth/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	err string
}

func (e ErrorResponse) Error() string {
	return e.err
}

func extractTokenIDFromHeader(w http.ResponseWriter, r * http.Request) (*model.Token, error){
	header := r.Header.Get("x-access-token")
	header = strings.TrimSpace(header)
	if header == "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string {"message":"Missing Auth token"})
	}

	token, err := middleware.ValidateToken(header)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "Error retrieving token",
		})
	}
	return token, nil
}

func GetPostsByPriority(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)
	if err != nil {
		return
	}
	//So, idk how to do this point. I think that is should be a GET endpoint but then it means that I don't have
	// parameters unless they are in the query. That means to send it in the query string and then use token +
	// the query string to build the query to the database. Check if that makes sense
	queryValues := r.URL.Query()
	if queryValues.Get("priority") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "Error retrieving token",
		})
	}
	datastore.GetPostByPriority(token.UserID, queryValues.Get("priority"))
}
func GetPosts(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)
	if err != nil {
		return
	}
	todoList := datastore.GetPosts(token.UserID)

	json.NewEncoder(w).Encode(todoList)
}

func AddPost(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)

	if err != nil {
		return
	}

	todo := &model.Todo{}
	json.NewDecoder(r.Body).Decode(todo)
	if todo.TimeToRead == 0 {
		todo.TimeToRead, err = utils.TimeToRead(todo.URL)
		if err != nil {
			todo.TimeToRead = 0
			fmt.Println("Error trying to calculate the time to read: ", err.Error())
		}
	}
	err = datastore.AddPost(token.UserID, todo)

}

func CreateUser(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	user := &model.User{}

	json.NewDecoder(r.Body).Decode(user)

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error, ", err.Error())
		err := ErrorResponse{
			err: "Password encryption failed",
		}
		json.NewEncoder(w).Encode(err)
	}
	user.Password = string(password)
	_, err = datastore.AddUser(user)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			err: "Error creating the user in the database",
		}
		json.NewEncoder(w).Encode(err)
	}
}

func Login(w http.ResponseWriter, r* http.Request, datastore datastore.Datastore) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		resp := map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := datastore.FindUser(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

func GetUsers(w http.ResponseWriter, r* http.Request, datastore datastore.Datastore) {
	users := datastore.FindUsers()
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r * http.Request, datastore datastore.Datastore) {
	params := mux.Vars(r)
	id := params["id"]
	user := datastore.FindUserById(id)
	json.NewEncoder(w).Encode(user)
}