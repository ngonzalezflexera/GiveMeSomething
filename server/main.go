package main

import (
	//"auth/datastore"
	"auth/routes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	e := godotenv.Load()

	if e != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(e)

	port := os.Getenv("PORT")

	//db := datastore.NewDatastore("postgres")
	//defer db.Close()
	// Handle routes
	http.Handle("/", routes.Handlers(nil))

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
