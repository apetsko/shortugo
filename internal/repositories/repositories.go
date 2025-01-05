package repositories

type Storage interface {
	Put(string) (string, error)
	Get(string) (string, error)
}
