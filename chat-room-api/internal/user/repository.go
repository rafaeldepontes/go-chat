package user

import "github.com/rafaeldepontes/go-chat/internal/model"

type Repository interface {
	FindAll() ([]model.User, error)
	Save(m *model.User) error
}
