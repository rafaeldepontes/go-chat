package user

import (
	"github.com/rafaeldepontes/go-chat/internal/model"
	pb "github.com/rafaeldepontes/go-chat/shared/message"
)

type Repository interface {
	FindAll() ([]*pb.Message, error)
	Save(m *model.User) error
}
