package dataaccess

// CRUD for the topics
// func CreateTopicTables(db *Database ) (error) {
// 	// Create tables if it doesn't exist
//     createTopicTableSQL := `

// 		CREATE TABLE IF NOT EXISTS topics (
//             "topic_id" SERIAL PRIMARY KEY NOT NULL,
//             "topic_name" TEXT NOT NULL,
// 			"admin_id" INT,
// 			"creator_id" INT NOT NULL,
// 			"date_created" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

// 			CONSTRAINT fk_admin FOREIGN KEY (admin_id)
// 			REFERENCES topics(topic_id),
// 			CONSTRAINT fk_topic_creator FOREIGN KEY (creator_id)
// 			REFERENCES topics(topic_id) ON DELETE SET NULL
//         );
//     `

//     _, err := db.Conn.Exec(context.Background(), createTopicTableSQL)
//     if err != nil {
//         return err
//     } else {
// 		fmt.Println("topics table created successfully")
// 		return nil
// 	}
// }

import (
	"context"

	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/models"
)

// type Topic struct {
// 	TopicID 	int			`json:"topicid"`
// 	TopicName 	string 		`json:"topicname"`
// 	AdminID   	int    		`json:"adminid"`
// 	CreatorID 	int 		`json:"creatorid"`
// 	DateCreated time.Time	`json:"date_created"`
// }

// Reading
func GetAllTopics(db *database.Database) ([]*models.Topic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT topic_id, topic_name, admin_id, creator_id, date_created FROM topics`
	rows, err := db.Conn.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.TopicID,
			&topic.TopicName,
			&topic.AdminID,
			&topic.CreatorID,
			&topic.DateCreated,
		)

		if err != nil {
			return nil, err
		}

		topics = append(topics, &topic)

	}
	// remember that it returns a pointer to the slice, and not the slice itself
	return topics, nil
}

// func GetTopicByName

func GetTopicByID(db *database.Database, id int) (*models.Topic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT topic_id, topic_name, admin_id, creator_id, date_created FROM topics WHERE topic_id = $1`
	row := db.Conn.QueryRow(ctx, query, id)

	var topic models.Topic

	err := row.Scan(
		&topic.TopicID,
		&topic.TopicName,
		&topic.AdminID,
		&topic.CreatorID,
		&topic.DateCreated,
	)

	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// Writing

func CreateTopic(db *database.Database, topic models.Topic) (*models.Topic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO topics (topic_name, admin_id, creator_id)
	VALUES ($1, $2, $3) returning topic_id, topic_name, admin_id, creator_id, date_created`

	var returnedTopic models.Topic

	err := db.Conn.QueryRow(
		ctx,
		query,
		topic.TopicName,
		topic.AdminID,
		topic.CreatorID,
	).Scan(
		&returnedTopic.TopicID,
		&returnedTopic.TopicName,
		&returnedTopic.AdminID,
		&returnedTopic.CreatorID,
		&returnedTopic.DateCreated,
	)
	if err != nil {
		return nil, TranslateTopicError(err)
	}

	return &returnedTopic, nil
}

func UpdateTopic(db *database.Database, id int, body models.Topic) (*models.Topic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	query := `
        UPDATE topics
        SET
            topic_name = $1, 
            admin_id = $2
        WHERE topic_id = $3
        returning topic_id, topic_name, admin_id, creator_id, date_created
		`
	var returnedtopic models.Topic

	err := db.Conn.QueryRow(
		ctx,
		query,
		body.TopicName,
		body.AdminID,
		id,
	).Scan(
		&returnedtopic.TopicID,
		&returnedtopic.TopicName,
		&returnedtopic.AdminID,
		&returnedtopic.CreatorID,
		&returnedtopic.DateCreated,
	)
	if err != nil {
		return nil, TranslateTopicError(err)
	}

	return &returnedtopic, nil

}

func DeleteTopic(db *database.Database, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM topics WHERE topic_id = $1`
	_, err := db.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
