package reader

import "io/ioutil"

type File struct{}

func (f *File) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
