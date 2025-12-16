package service

import (
	"context"
	"encoding/json"

	msgInterface "github.com/rafaeldepontes/go-chat/internal/message"
	"github.com/rafaeldepontes/go-chat/internal/model"
	messageapi "github.com/rafaeldepontes/go-chat/pkg/gRPC/message-api"
	pb "github.com/rafaeldepontes/go-chat/shared/message"
)

type messageSvc struct {
	msgClient *pb.MessageServiceClient
}

func NewService() msgInterface.Service {
	conn := messageapi.GetgRPCServer()
	client := pb.NewMessageServiceClient(conn)
	return &messageSvc{
		msgClient: &client,
	}
}

func (m *messageSvc) FindAll() ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gRPCResp, err := (*m.msgClient).FindAll(ctx, &pb.MessageRequest{})
	if err != nil {
		return nil, err
	}

	msgRPCResp := gRPCResp.Data
	msgs := make([]model.User, 0, len(msgRPCResp))

	for _, val := range msgRPCResp {
		msg := model.User{
			Username: val.Username,
			Message:  val.Message,
		}
		msgs = append(msgs, msg)
	}

	parsedJson, err := json.Marshal(msgs)
	if err != nil {
		return nil, err
	}

	return parsedJson, nil
}
