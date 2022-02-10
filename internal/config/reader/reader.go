package reader

type Reader interface {
	Read(string) ([]byte, error)
}
