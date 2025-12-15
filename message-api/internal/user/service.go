package user

type Service interface {
	GetUserChannel() *chan []byte
	FindAll() ([]byte, error)
	Save()
}
