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
	
	// row := conn.QueryRow(context.Background(), "SELECT current_database(), current_schema()")
	// var dbName, schema string
	// row.Scan(&dbName, &schema)
	// fmt.Println("Connected to DB:", dbName, "Schema:", schema)

	return &Database{Conn: conn}, nil
}

// DROP TABLE IF EXISTS comments;
// DROP TABLE IF EXISTS posts;
// DROP TABLE IF EXISTS topics;
// DROP TABLE IF EXISTS users;


func CreateUserTables(db *Database ) (error) {
	// Create tables if it doesn't exist
    createUserTableSQL := `

		CREATE TABLE IF NOT EXISTS users (
            "user_id" SERIAL PRIMARY KEY NOT NULL,
            "username" varchar(255) UNIQUE NOT NULL,
			"password" varchar(255) NOT NULL,
			"role" varchar(255) NOT NULL DEFAULT 'User',
            "date_created" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
    `

    _, err := db.Conn.Exec(context.Background(), createUserTableSQL)
    if err != nil {
        return err
    } else {
		fmt.Println("users table created successfully")
		return nil
	}
}


func CreateTopicTables(db *Database ) (error) {
	// Create tables if it doesn't exist
    createTopicTableSQL := `

		CREATE TABLE IF NOT EXISTS topics (
            "topic_id" SERIAL PRIMARY KEY NOT NULL,
            "topic_name" TEXT NOT NULL,
			"admin_id" INT,
			"creator_id" INT NOT NULL,
			"date_created" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

			CONSTRAINT fk_admin FOREIGN KEY (admin_id) 
			REFERENCES users(user_id),
			CONSTRAINT fk_topic_creator FOREIGN KEY (creator_id) 
			REFERENCES users(user_id) ON DELETE SET NULL  
        );
    `

    _, err := db.Conn.Exec(context.Background(), createTopicTableSQL)
    if err != nil {
        return err
    } else {
		fmt.Println("topics table created successfully")
		return nil
	}
}


func CreatePostTables(db *Database ) (error) {
	// Create tables if it doesn't exist
    createPostTableSQL := `
		CREATE TABLE IF NOT EXISTS posts (
            "post_id" SERIAL PRIMARY KEY NOT NULL,
            "Title" TEXT NOT NULL,
			"content" TEXT NOT NULL,
			"creator_id" INT NOT NULL,
			"topic_id" INT NOT NULL,
            "date_created" TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
			"date_updated" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),


			CONSTRAINT fk_creator FOREIGN KEY (creator_id) 
			REFERENCES users(user_id) ON DELETE CASCADE, 
			CONSTRAINT fk_topic FOREIGN KEY (topic_id) 
			REFERENCES topics(topic_id) ON DELETE CASCADE
        );
    `

    _, err := db.Conn.Exec(context.Background(), createPostTableSQL)
    if err != nil {
        return err
    } else {
		fmt.Println("posts table created successfully")
		return nil
	}

}


func CreateCommentTables(db *Database ) (error) {
	// Create tables if it doesn't exist
    createCommentTableSQL := `
		CREATE TABLE IF NOT EXISTS comments (
            "comment_id" SERIAL PRIMARY KEY NOT NULL,
            "content" TEXT NOT NULL,
			"post_id" INT NOT NULL,
			"creator_id" INT NOT NULL,
			"topic_id" INT NOT NULL,
			"parent_comment_id" INT DEFAULT NULL, 
            "date_created" TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
			"date_updated" TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
			
			CONSTRAINT fk_post FOREIGN KEY (post_id) 
			REFERENCES posts(post_id) ON DELETE CASCADE,
			CONSTRAINT fk_creator FOREIGN KEY (creator_id) 
			REFERENCES users(user_id) ON DELETE CASCADE, 
			CONSTRAINT fk_topic FOREIGN KEY (topic_id) 
			REFERENCES topics(topic_id) ON DELETE CASCADE,
			CONSTRAINT fk_parent_comment FOREIGN KEY (parent_comment_id) 
			REFERENCES comments(comment_id) ON DELETE CASCADE
        );
    `

    _, err := db.Conn.Exec(context.Background(), createCommentTableSQL)
    if err != nil {
        return err
    } else {
		fmt.Println("comments table created successfully")
		return nil
	}

}

