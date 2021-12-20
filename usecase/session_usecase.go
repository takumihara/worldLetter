package usecase

import (
	"github.com/tacomea/worldLetter/domain"
	"github.com/tacomea/worldLetter/domain/repository"
)

type SessionUseCase interface {
	Create(session domain.Session) error
	Delete(id string) error
	Read(id string) (domain.Session, error)
}


type sessionUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewSessionUsecase(sr repository.SessionRepository) SessionUseCase {
	return &sessionUsecase{
		sessionRepo: sr,
	}
}

func (u *sessionUsecase) Create(session domain.Session) error {
	return u.sessionRepo.Create(session)
}

func (u *sessionUsecase) Delete(id string) error {
	return u.sessionRepo.Delete(id)
}

func (u *sessionUsecase) Read(id string) (domain.Session, error) {
	return u.sessionRepo.Read(id)
}
