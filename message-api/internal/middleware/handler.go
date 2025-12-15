package middleware

import (
	"net/http"

	messagebroker "github.com/rafaeldepontes/go-chat/internal/message-broker"
	msgBrokerSvc "github.com/rafaeldepontes/go-chat/internal/message-broker/service"
	user "github.com/rafaeldepontes/go-chat/internal/user"
	userSvc "github.com/rafaeldepontes/go-chat/internal/user/service"
)

type handler struct {
	MessageBroker messagebroker.MsgBroker
	UserService   user.Service
}

func NewHandler() *handler {
	return &handler{
		UserService:   userSvc.NewService(),
		MessageBroker: msgBrokerSvc.NewMsgService(),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}
