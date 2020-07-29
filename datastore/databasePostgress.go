package datastore

import (
	"auth/model"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres dialect interface
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"time"

	"log"
	"os"
)

type postgresDatastore struct {
	client * gorm.DB
}

func (p *postgresDatastore) FindUserById(id string) *model.User {
	user := &model.User{}
	p.client.First(user, id)
	return user
}

func (p *postgresDatastore) GetPostByPriority(id uint, priority string) []*model.Todo {
	userModel := &model.User{}
	if err := p.client.Where("id=?",id).Preload("Todos", "priority = (?)", priority).
		Find(userModel).Error; err != nil {
			fmt.Println("Error retrieving the todos for the user ", userModel.ID)
			return []*model.Todo{}
	}
	return userModel.Todos
}

func (p *postgresDatastore) FindUsers() [] model.User {
	var users [] model.User
	p.client.Preload("auths").Find(&users)
	return users
}

func (p postgresDatastore) Close() {
	p.client.Close()
}

func (p *postgresDatastore) AddPost(id uint, todo *model.Todo) error {
	todo.UserID = id
	if err := p.client.Create(todo).Error; err != nil {
		return errors.New("error adding item")
	}
	return nil
}

func (p *postgresDatastore) GetPosts(id uint) []*model.Todo {
	var ret []*model.Todo
	user := &model.User{
		Model: gorm.Model{
			ID: id,
		},
	}
	if err := p.client.Model(user).Related(&ret).Error; err != nil {
		return []*model.Todo{}
	}
	return ret
}


func createPostgresDatastore() Datastore {
	datastore := &postgresDatastore{}
	client, err :=  connectDB()
	if err != nil {
		log.Fatal("Error creating connection to the postgres database", err.Error())
	}
	datastore.client = client
	return datastore
}
func ( p * postgresDatastore) AddUser(user *model.User) (status bool, err error){
	result := p.client.Create(user)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (p postgresDatastore) FindUser(email, password string) map[string]interface{} {
	user := &model.User{}
	if err := p.client.Where("Email = ?", email).First(user).Error; err !=nil {
		resp := map[string]interface{}{"status":false, "message":"Email address not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute* 100000).Unix()
	err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		resp := map[string]interface{}{"status":false, "message":"Password doesn't match"}
		return resp
	}

	token := &model.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	jwtoken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenString, err := jwtoken.SignedString([]byte("secretword"))

	if err != nil {
		fmt.Println(err)
	}

	resp := map[string]interface{} {"status": false, "message": "Logged in"}
	resp["token"] = tokenString
	resp["user"] = user
	return resp
}


func connectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading the .env file")
	}
	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	dbString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",databaseHost, username, databaseName, password)
	db, err := gorm.Open("postgres", dbString)
	db.LogMode(true)

	if err != nil {
		return nil, errors.New("error connecting to the database, " + err.Error())
	}
	db.AutoMigrate(
		&model.User{},
		&model.Todo{})
	fmt.Print("Successful connection")
	db.Model(&model.Todo{}).AddForeignKey("user_id", "Users(id)", "CASCADE", "RESTRICT")

	return db, nil
}

func mockData (db * gorm.DB) {
	//user := &model.User{
	//	Name:     "Test1",
	//	Email:    "test@test.com",
	//	Password: "test",
	//}
	//user2 := &model.User{
	//	Name:     "Test2",
	//	Email:    "Test2@test.com",
	//	Password: "test2",
	//}
	//todo1:= &model.Todo{
	//	UserID:      10,
	//	Title:       "Test",
	//	URL:         "www.test.com",
	//	Description: "testtest",
	//	TimeToRead:  3,
	//	Priority:    1,
	//}
	//todo2:= &model.Todo{
	//	UserID:      10,
	//	Title:       "Test2",
	//	URL:         "www.test2.com",
	//	Description: "testtest2",
	//	TimeToRead:  3,
	//	Priority:    1,
	//}
	//todo3:= &model.Todo{
	//	UserID:      10,
	//	Title:       "Test3",
	//	URL:         "www.test3.com",
	//	Description: "testtest3",
	//	TimeToRead:  3,
	//	Priority:    1,
	//}
	//todo4:= &model.Todo{
	//	Title:       "Test4",
	//	URL:         "www.test.com4",
	//	Description: "testtest4",
	//	TimeToRead:  3,
	//	Priority:    1,
	//}
	//res := db.Create(user)
	//fmt.Println(res.Error, " ", res.Value , " ", res.RowsAffected)
	//user2.Todos = append(user2.Todos, todo4)
	//res = db.Create(user2)
	//fmt.Println(res.Error, " ", res.Value , " ", res.RowsAffected)
	//res = db.Create(todo1)
	//fmt.Println(res.Error, " ", res.Value , " ", res.RowsAffected)
	//res = db.Create(todo2)
	//fmt.Println(res.Error, " ", res.Value , " ", res.RowsAffected)
	//res = db.Create(todo3)
	//fmt.Println(res.Error, " ", res.Value , " ", res.RowsAffected)


	//// Method 1
	//bar := &[]model.Todo{}
	//finduser := model.User{
	//	Model: gorm.Model{
	//		ID:        76,
	//	},
	//}
	//db.Debug().Model(&finduser).Related(bar)
	////Method 2
	foo := &model.User{}
	db.Debug().Where("id=?",76).Preload("Todos", "priority = (?)",1).Find(foo)


	//Method 3
	//rows, err := db.Table("users").Where("users.id = ?", 76).Joins(
	//	"Join todos on todos.user_id = users.id").Where("todos.priority = ?", 1).Select(
	//		"users.id, users.name, todos.url").Rows()
	//if err != nil {
	//	log.Panic(err)
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var fii string
	//	var faa string
	//	var fuu string
	//	err := rows.Scan(&fii, &faa, &fuu)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//}
}