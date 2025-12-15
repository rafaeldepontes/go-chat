package user

type Service interface {
	FindAll() ([]byte, error)
	Save(msg ...byte) error
}
