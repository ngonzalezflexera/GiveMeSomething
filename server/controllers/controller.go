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
	Err string
}

func (e ErrorResponse) Error() string {
	return e.Err
}

// extractTokenIDFromHeader will extract the token from the header and it will validate it.
// If the token is empty, it will return a status forbidden..
// If the token is invalid, it will return a bad request.
// On a valid token request, it will return the model of a token
func extractTokenIDFromHeader(w http.ResponseWriter, r *http.Request) (*model.Token, error) {
	header := r.Header.Get("x-access-token")
	header = strings.TrimSpace(header)
	if header == "" {
		w.WriteHeader(http.StatusForbidden)
		err := json.NewEncoder(w).Encode(map[string]string{"message": "Missing Auth token"})
		if err != nil {
			fmt.Println("Error encoding the response")
		}
	}

	token, err := middleware.ValidateToken(header)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Error retrieving token",
		})
		if err != nil {
			fmt.Println("Error encoding the response")
		}
	}
	return token, nil
}

// GetPostsByPriority will return all the post for an user, with a specified priority. The priority will be found
// in the query string. The user will be embed in the token.
// On valid request it will return all the todos for a specific priority
func GetPostsByPriority(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)
	if err != nil {
		return
	}
	queryValues := r.URL.Query()
	if queryValues.Get("priority") == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Error retrieving token",
		})
		if err != nil {
			fmt.Println("Error encoding the response")
		}
	}
	todolist := datastore.GetPostByPriority(token.UserID, queryValues.Get("priority"))
	err = json.NewEncoder(w).Encode(todolist)
	if err != nil {
		fmt.Println("Error encoding the response")
	}
}

// GetPosts will return all the post for an user. The user will be embed in the token
// On valid request it will return all the todos for a specific priority
func GetPosts(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)
	if err != nil {
		return
	}
	todoList := datastore.GetPosts(token.UserID)

	err = json.NewEncoder(w).Encode(todoList)
	if err != nil {
		fmt.Println("Error encoding the response")
	}
}

// AddPost will add a new post to the user.
// The user will be embed in the token.
// On valid request it will return a status 200
func AddPost(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	token, err := extractTokenIDFromHeader(w, r)

	if err != nil {
		return
	}

	todo := &model.Todo{}
	err = json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		fmt.Println("Error trying to decode the TODO")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if todo.TimeToRead == 0 {
		todo.TimeToRead, err = utils.TimeToRead(todo.URL)
		if err != nil {
			todo.TimeToRead = 0
			fmt.Println("Error trying to calculate the time to read: ", err.Error())
		}
	}
	err = datastore.AddPost(token.UserID, todo)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// CreateUser will create a user in the datastore. The password will be encrypted and saved in the datastore
func CreateUser(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("Error trying to decode the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error, ", err.Error())
		err := ErrorResponse{
			Err: "Password encryption failed",
		}
		w.WriteHeader(http.StatusInternalServerError)
		err2 := json.NewEncoder(w).Encode(err)
		if err2 != nil {
			fmt.Println("Error encoding the response")
		}
	}
	user.Password = string(password)
	err = datastore.AddUser(user)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Error creating the user in the database",
		}
		err2 := json.NewEncoder(w).Encode(err)
		if err2 != nil {
			fmt.Println("Error encoding the response")
		}
	}
	w.WriteHeader(http.StatusOK)
}

// Login will establish if the user is a valid user.
func Login(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		resp := map[string]interface{}{"status": false, "message": "Invalid request"}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			fmt.Println("Error encoding the response")
		}
		return
	}
	resp, err := datastore.FindUser(user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			fmt.Println("Error encoding the response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error encoding the response")
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	users := datastore.FindUsers()
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request, datastore datastore.Datastore) {
	params := mux.Vars(r)
	id := params["id"]
	user := datastore.FindUserById(id)
	json.NewEncoder(w).Encode(user)
}
