package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"log"
	"sync"
)

type letterRepository struct {
	m sync.Map
}

func NewSyncMapLetterRepository() domain.LetterRepository {
	return &letterRepository{}
}

func (s *letterRepository) Create(letter domain.Letter) error {
	s.m.Store(letter.ID, letter)
	return nil
}

func (s *letterRepository) Read(id string) (domain.Letter, error) {
	if value, ok := s.m.Load(id); ok {
		return value.(domain.Letter), nil
	}
	return domain.Letter{}, errors.New("user not found")
}

func (s *letterRepository) Update(letter domain.Letter) error {
	s.m.Store(letter.ID, letter)
	_, ok := s.m.Load(letter.ID)
	if !ok {
		return errors.New("user not created")
	}
	return nil
}

func (s *letterRepository) Delete(id string) error {
	s.m.Delete(id)
	return nil
}

func (s *letterRepository) 	ShowUnsendRandomLetter(AuthorID string) (domain.Letter, error) {
	var letter domain.Letter
	s.m.Range(func(key interface{}, value interface{}) bool {
		val := value.(domain.Letter)
		if val.IsSent == false && val.AuthorID != AuthorID{
			letter = val
			log.Println("letter was found!")
			return false
		}
		return true
	})
	log.Println(letter)
	return letter, nil
}
