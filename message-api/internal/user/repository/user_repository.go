package repository

import (
	"database/sql"
	"fmt"

	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/user"
	"github.com/rafaeldepontes/go-chat/pkg/db/postgres"

	pb "github.com/rafaeldepontes/go-chat/shared/message"
)

type userRepo struct {
	db *sql.DB
}

func NewRepository() user.Repository {
	return &userRepo{
		db: postgres.GetDb(),
	}
}

func (repo *userRepo) FindAll() ([]*pb.Message, error) {
	fmt.Println("Listing all the messages...")
	var msgs []*pb.Message

	var rows *sql.Rows
	rows, err := repo.db.Query(`SELECT username, message FROM chat_room cr;`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg pb.Message
		if err = rows.Scan(&msg.Username, &msg.Message); err != nil {
			return nil, err
		}
		msgs = append(msgs, &msg)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("Found %v messages\n", len(msgs))
	return msgs, nil
}

func (repo *userRepo) Save(user *model.User) error {
	fmt.Println("Saving message:", user.Message)
	_, err := repo.db.Exec(`INSERT INTO chat_room (username, message) VALUES ($1, $2)`, user.Username, user.Message)
	return err
}
