package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/gorm"
)

type letterRepositoryPG struct {
	db *gorm.DB
}

func NewLetterRepositoryPG(db *gorm.DB) domain.LetterRepository {
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

func (l *letterRepositoryPG) GetFirstUnsendLetter(AuthorID string) (domain.Letter, error) {
	var letter domain.Letter
	l.db.Find(&letter, "is_sent = ? AND author_id != ?", false, AuthorID)
	return letter, nil
}