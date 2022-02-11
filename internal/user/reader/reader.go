package reader

type Config interface {
	Read(string) ([]byte, error)
}
