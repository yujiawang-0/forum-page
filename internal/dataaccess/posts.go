package dataaccess

import (
	"context"

	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/models"
)

// Reading
func GetAllPosts(db *database.Database) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT post_id, title, content, creator_id, topic_id, date_created, date_updated FROM posts`
	rows, err := db.Conn.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.Title,
			&post.Content,
			&post.CreatorID,
			&post.TopicID,
			&post.DateCreated,
			&post.DateUpdated,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)

	}
	// remember that it returns a pointer to the slice, and not the slice itself
	return posts, nil
}

func GetPostsByTopicId(db *database.Database, topicid int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT post_id, title, content, creator_id, topic_id, date_created, date_updated 
	FROM posts 
	WHERE topic_id = $1
	ORDER BY date_created DESC;
	`
	rows, err := db.Conn.Query(ctx, query, topicid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.Title,
			&post.Content,
			&post.CreatorID,
			&post.TopicID,
			&post.DateCreated,
			&post.DateUpdated,
		)
		if err != nil {
		return nil, err
		}
		
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(db *database.Database, id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT post_id, title, content, creator_id, topic_id, date_created, date_updated FROM posts WHERE post_id = $1`
	row := db.Conn.QueryRow(ctx, query, id)

	var post models.Post

	err := row.Scan(
		&post.PostID,
		&post.Title,
		&post.Content,
		&post.CreatorID,
		&post.TopicID,
		&post.DateCreated,
		&post.DateUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Writing

func CreatePost(db *database.Database, post models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO posts (title, content, creator_id, topic_id)
	VALUES ($1, $2, $3, $4) returning post_id, title, content, creator_id, topic_id, date_created, date_updated`

	var new models.Post

	err := db.Conn.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.CreatorID,
		post.TopicID,
	).Scan(
		&new.PostID,
		&new.Title,
		&new.Content,
		&new.CreatorID,
		&new.TopicID,
		&new.DateCreated,
		&new.DateUpdated,
	)
	if err != nil {
		return nil, TranslatePostError(err)
	}

	return &new, nil
}

func UpdatePost(db *database.Database, id int, body models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        UPDATE posts
        SET
            title = $1, 
            content = $2,
			date_updated = NOW()
        WHERE post_id = $3
        returning post_id, title, content, creator_id, topic_id, date_created, date_updated
		`
	var updated models.Post

	err := db.Conn.QueryRow(
		ctx,
		query,
		body.Title,
		body.Content,
		id,
	).Scan(
		&updated.PostID,
		&updated.Title,
		&updated.Content,
		&updated.CreatorID,
		&updated.TopicID,
		&updated.DateCreated,
		&updated.DateUpdated,
	)
	if err != nil {
		return nil, TranslatePostError(err)
	}

	return &updated, nil

}

func DeletePost(db *database.Database, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM posts WHERE post_id = $1`
	_, err := db.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
