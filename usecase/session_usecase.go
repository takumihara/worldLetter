package usecase

import "github.com/tacomea/worldLetter/domain"

type sessionUsecase struct {
	sessionRepo domain.SessionRepository
}

func NewSessionUsecase(sr domain.SessionRepository) domain.SessionUseCase {
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
