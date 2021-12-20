package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/gorm"
)

type letterRepositoryPG struct {
	db *gorm.DB
}

func NewLetterRepositoryPG(db *gorm.DB) LetterRepository {
	return &letterRepositoryPG{
		db: db,
	}
}

func (l *letterRepositoryPG) Create(letter domain.Letter) error {
	result := l.db.Create(letter)
	if result.Error != nil {
		return errors.New("unexpected error while storing letter")
	}
	return nil
}

func (l *letterRepositoryPG) Update(letter domain.Letter) error {
	l.db.Save(&letter)
	return nil
}

func (l *letterRepositoryPG) Read(id string) (domain.Letter, error) {
	var letter domain.Letter
	result := l.db.First(&letter, "id = ?", id)
	if result.Error != nil {
		return domain.Letter{}, errors.New("unexpected error while retrieving user")
	}
	return letter, nil
}

func (l *letterRepositoryPG) Delete(id string) error {
	result := l.db.Delete(&domain.Letter{}, id)
	if result.Error != nil {
		return errors.New("unexpected error while deleting letter")
	}
	return nil
}

func (l *letterRepositoryPG) GetAll() ([]domain.Letter, error) {
	var letters []domain.Letter

	results := l.db.Find(&letters)
	if results.Error != nil {
		return nil, results.Error
	}

	return letters, nil
}

func (l *letterRepositoryPG) GetFirstUnsendLetter(authorID string) (domain.Letter, error) {
	var letter domain.Letter
	l.db.Find(&letter, "is_sent = ? AND author_id != ?", false, authorID)
	return letter, nil
}

func (l *letterRepositoryPG) GetLettersByAuthorID(authorID string) (string, error) {
	var str string
	var letter domain.Letter
	l.db.Find(&letter, "is_sent = ? AND author_id != ?", false, authorID)
	return str, nil
}

func (l *letterRepositoryPG) GetLettersByReceiverID(receiverID string) (string, error) {
	var str string
	var letter domain.Letter
	l.db.Find(&letter, "is_sent = ? AND receiver_id != ?", false, receiverID)
	return str, nil
}