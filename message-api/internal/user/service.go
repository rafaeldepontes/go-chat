package user

import (
	"context"

	pb "github.com/rafaeldepontes/go-chat/shared/message"
)

type Service interface {
	GetUserChannel() *chan []byte
	FindAll(context.Context, *pb.MessageRequest) (*pb.MessageResponses, error)
	Save()
	Run() error
}
