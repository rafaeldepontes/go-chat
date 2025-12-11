package service

import (
	"encoding/json"
	"time"

	"github.com/rafaeldepontes/go-chat/internal/cache"
	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/user"
	"github.com/rafaeldepontes/go-chat/internal/user/repository"
)

type userSvc struct {
	cache      *cache.Cache[string, []model.User]
	repository user.Repository
}

func NewService() user.Service {
	return &userSvc{
		repository: repository.NewRepository(),
		cache:      cache.NewCache[string, []model.User](),
	}
}

func (s *userSvc) FindAll() ([]byte, error) {
	val, ok := s.cache.Get("messages")
	if ok {
		return json.Marshal(*val)
	}

	chats, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	s.cache.Set("messages", chats, time.Duration(1))

	return json.Marshal(chats)
}

func (s *userSvc) Save(m ...byte) error {
	var user model.User
	if err := json.Unmarshal(m, &user); err != nil {
		return err
	}
	return s.repository.Save(&user)
}
