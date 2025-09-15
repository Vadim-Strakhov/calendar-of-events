package storage

type Store interface {
	Save(data []byte) error
	Load() ([]byte, error)
	GetFilename() string
}

type Storage struct {
	filename string
}

func NewStorage(filename string) *Storage {
	return &Storage{filename: filename}
}

func (s *Storage) GetFilename() string {
	return s.filename
}
