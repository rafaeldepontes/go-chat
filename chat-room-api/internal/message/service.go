package message

type Service interface {
	FindAll() ([]byte, error)
}
