package repository

import "github.com/tacomea/worldLetter/domain"

type LetterRepository interface {
	Create(letter domain.Letter) error
	Read(id string) (domain.Letter, error)
	Update(letter domain.Letter) error
	Delete(id string) error
	GetAll() ([]domain.Letter, error)
	GetFirstUnsendLetter(authorID string) (domain.Letter, error)
	GetLettersByAuthorID(authorID string) (string, error)
	GetLettersByReceiverID(receiverID string) (string, error)
}

