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
	err := u.sessionRepo.Create(session)
	if err != nil {
		return err
	}
	return nil
}

func (u *sessionUsecase) Delete(id string) error {
	err := u.sessionRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *sessionUsecase) Read(id string) (domain.Session, error) {
	user, err := u.sessionRepo.Read(id)
	if err != nil {
		return domain.Session{}, err
	}
	return user, nil
}
