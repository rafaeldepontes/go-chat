package messagebroker

type MsgBroker interface {
	Send(msg []byte) error
}
