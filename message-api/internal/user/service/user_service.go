package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/rafaeldepontes/go-chat/internal/model"
	"github.com/rafaeldepontes/go-chat/internal/user"
	"github.com/rafaeldepontes/go-chat/internal/user/repository"
	pb "github.com/rafaeldepontes/go-chat/shared/message"
	"google.golang.org/grpc"
)

type userSvc struct {
	pb.UnimplementedMessageServiceServer
	grpcServer *grpc.Server
	repository user.Repository
	msg        chan []byte
}

func NewService() user.Service {
	grpcServer := grpc.NewServer()

	userSvc := &userSvc{
		grpcServer: grpcServer,
		repository: repository.NewRepository(),
		msg:        make(chan []byte),
	}

	pb.RegisterMessageServiceServer(grpcServer, userSvc)

	return userSvc
}

func (s *userSvc) GetUserChannel() *chan []byte {
	return &s.msg
}

func (s *userSvc) FindAll(context.Context, *pb.MessageRequest) (*pb.MessageResponses, error) {
	chats, err := s.repository.FindAll()

	if err != nil {
		return nil, err
	}

	resp := &pb.MessageResponses{
		Data: chats,
	}

	return resp, nil
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

func (s *userSvc) Run() error {
	var port string
	port = ":" + os.Getenv("SERVER_PORT")

	list, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("[Error] connection to: ", os.Getenv("SERVER_PORT"))
	}

	return s.grpcServer.Serve(list)
}
