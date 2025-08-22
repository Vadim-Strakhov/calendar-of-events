package storage

// Store описывает контракт для любых стораджей
type Store interface {
	Save(data []byte) error
	Load() ([]byte, error)
	GetFilename() string
}

// Storage — базовый тип с общими полями/методами
type Storage struct {
	filename string
}

// NewStorage создает базовый Storage (может использоваться внутри конкретных стораджей)
func NewStorage(filename string) *Storage {
	return &Storage{filename: filename}
}

// GetFilename возвращает имя файла стораджа
func (s *Storage) GetFilename() string {
	return s.filename
}
