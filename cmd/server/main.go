package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/router"
)

func main() {
	db, err:= database.GetDB()
	if err!= nil {
		log.Fatalf("Unable to get DB: %v", err)
	}

	defer db.Conn.Close(context.Background())
 
	// Create table if it doesn't exist
    createTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err = db.Conn.Exec(context.Background(), createTableSQL)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create table: %v\n", err)
        os.Exit(1)
    } else {
		fmt.Println("users table created successfully")
	}

	r := router.Setup()
	fmt.Print("Listening on port 8000 at http://localhost:8000!")

	log.Fatalln(http.ListenAndServe(":8000", r))
}
