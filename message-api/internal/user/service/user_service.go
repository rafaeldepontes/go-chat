package service

import (
	"encoding/json"
	"fmt"

	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/user"
	"github.com/rafaeldepontes/go-chat/internal/user/repository"
)

type userSvc struct {
	repository user.Repository
	msg        chan []byte
}

func NewService() user.Service {
	return &userSvc{
		repository: repository.NewRepository(),
		msg:        make(chan []byte),
	}
}

func (s *userSvc) GetUserChannel() *chan []byte {
	return &s.msg
}

func (s *userSvc) FindAll() ([]byte, error) {
	chats, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	return json.Marshal(chats)
}

func (s *userSvc) Save() {
	var user model.User
	for data := range s.msg {
		if err := json.Unmarshal(data, &user); err != nil {
			fmt.Println("[ERROR] occurred trying to save a message:", err)
		}
		s.repository.Save(&user)
	}
}
