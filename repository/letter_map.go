package repository

import (
	"encoding/base64"
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"sync"
)

type letterRepository struct {
	m sync.Map
}

func NewSyncMapLetterRepository() domain.LetterRepository {
	return &letterRepository{}
}

func (l *letterRepository) Create(letter domain.Letter) error {
	l.m.Store(letter.ID, letter)
	return nil
}

func (l *letterRepository) Read(id string) (domain.Letter, error) {
	if value, ok := l.m.Load(id); ok {
		return value.(domain.Letter), nil
	}
	return domain.Letter{}, errors.New("user not found")
}

func (l *letterRepository) Update(letter domain.Letter) error {
	l.m.Store(letter.ID, letter)
	_, ok := l.m.Load(letter.ID)
	if !ok {
		return errors.New("user not created")
	}
	return nil
}

func (l *letterRepository) Delete(id string) error {
	l.m.Delete(id)
	return nil
}

func (l *letterRepository) GetFirstUnsendLetter(authorID string) (domain.Letter, error) {
	var letter domain.Letter
	l.m.Range(func(key interface{}, value interface{}) bool {
		val := value.(domain.Letter)
		if val.IsSent == false && val.AuthorID != authorID{
			letter = val
			return false
		}
		return true
	})
	return letter, nil
}

func (l *letterRepository) GetAll() ([]domain.Letter, error) {
	var letters []domain.Letter

	l.m.Range(func(key interface{}, value interface{}) bool {
		letters = append(letters, value.(domain.Letter))
		return true
	})

	return letters, nil
}

func (l *letterRepository) GetLettersByAuthorID(authorID string) (string, error) {
	var str string
	l.m.Range(func(key interface{}, value interface{}) bool {
		letter := value.(domain.Letter)
		if letter.AuthorID == authorID {
			encodedContent := base64.StdEncoding.EncodeToString([]byte(letter.Content))
			str += encodedContent + "|"
		}
		return true
	})
	return str, nil
}

func (l *letterRepository) GetLettersByReceiverID(receiverID string) (string, error) {
	var str string
	l.m.Range(func(key interface{}, value interface{}) bool {
		letter := value.(domain.Letter)
		if letter.ReceiverID == receiverID {
			encodedContent := base64.StdEncoding.EncodeToString([]byte(letter.Content))
			str += encodedContent + "|"
		}
		return true
	})
	return str, nil
}