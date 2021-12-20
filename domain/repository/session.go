package repository

import "github.com/tacomea/worldLetter/domain"

type SessionRepository interface {
	Create(session domain.Session) error
	Delete(id string) error
	Read(id string) (domain.Session, error)
}

