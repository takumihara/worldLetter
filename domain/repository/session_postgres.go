package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/gorm"
)

type sessionRepositoryPG struct {
	db *gorm.DB
}

func NewSessionRepositoryPG(db *gorm.DB) SessionRepository {
	return &sessionRepositoryPG{
		db: db,
	}
}

func (s *sessionRepositoryPG) Create(session domain.Session) error {
	result := s.db.Create(session)
	if result.Error != nil {
		return errors.New("unexpected error while storing session")
	}
	return nil
}

func (s *sessionRepositoryPG) Delete(id string) error {
	result := s.db.Delete(&domain.Session{}, id)
	if result.Error != nil {
		return errors.New("unexpected error while deleting session")
	}
	return nil
}

func (s *sessionRepositoryPG) Read(id string) (domain.Session, error) {
	var session domain.Session
	result := s.db.First(&session, "id = ?", id)
	if result.Error != nil {
		return domain.Session{}, errors.New("unexpected error while retrieving user")
	}
	return session, nil
}
