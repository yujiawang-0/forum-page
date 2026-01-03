package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	//"os"

	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/router"
	"github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/models"
)

func main() {
	db, err:= database.GetDB()
	if err!= nil {
		log.Fatalf("Unable to get DB: %v", err)
	}

	defer db.Conn.Close(context.Background()) // postpone this until main() function is over
	
	// create the tables in the database
	err = database.CreateUserTables(db);
	if err != nil {
		log.Fatal(err)
	}
	err = database.CreateTopicTables(db);
	if err != nil {
		log.Fatal(err)
	}
	err = database.CreatePostTables(db);
	if err != nil {
		log.Fatal(err)
	}
	err = database.CreateCommentTables(db);
	if err != nil {
		log.Fatal(err)
	}

	// testing dataaccess
	u, err := dataaccess.CreateUser(db, models.User{
		Username: "Tom",
		Password: "secret",
		Role: "User",
	})
	log.Println(u, err)

	users, _ := dataaccess.GetAllUsers(db)
	for _, u := range users {
		log.Println(*u)
	}
	// log.Println(users)

	one, _ := dataaccess.GetUserByID(db, u.ID)
	log.Println(one)

	
	r := router.Setup()
	fmt.Print("Listening on port 8000 at http://localhost:8000!")

	log.Fatalln(http.ListenAndServe(":8000", r))
}
