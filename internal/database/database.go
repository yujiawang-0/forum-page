/* what this file/folder does:
DB connections (and migrations)
*/

package database

import (
	"context"
	// "database/sql"
	"fmt"
	"os"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)
type Database struct {
	Conn *pgx.Conn
}


func GetDB() (*Database, error) {
	//connect to dotenv
	err := godotenv.Load(".env")
	if err != nil{
		log.Printf("Error loading .env file: %s", err)
		return nil, err
	}
	
	//connect to database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	} 
	
	fmt.Println("Successfully connected to database")
	return &Database{Conn: conn}, nil
}
