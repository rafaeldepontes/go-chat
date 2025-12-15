package messagebroker

type MsgBroker interface {
	Process(userChannel *chan []byte) error
}
