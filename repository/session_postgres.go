package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/gorm"
)

type sessionRepositoryMySQL struct {
	db *gorm.DB
}

func NewSessionRepositoryMySQL(db *gorm.DB) domain.SessionRepository{
	return &sessionRepositoryMySQL{
		db: db,
	}
}

func (s *sessionRepositoryMySQL) Create(session domain.Session) error {
	result := s.db.Create(session)
	if result.Error != nil {
		return errors.New("unexpected error while storing session")
	}
	return nil
}

func (s *sessionRepositoryMySQL) Delete(id string) error {
	result := s.db.Delete(&domain.Session{}, id)
	if result.Error != nil {
		return errors.New("unexpected error while deleting session")
	}
	return nil
}

func (s *sessionRepositoryMySQL) Read(id string) (domain.Session, error) {
	var session domain.Session
	result := s.db.First(&session, "id = ?", id)
	if result.Error != nil {
		return domain.Session{}, errors.New("unexpected error while retrieving user")
	}
	return session, nil
}


