package usecase

import "github.com/tacomea/worldLetter/domain"

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUseCase {
	return &userUsecase{
		userRepo: ur,
	}
}

func (u *userUsecase) Create(user domain.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) Read(email string) (domain.User, error) {
	return u.userRepo.Read(email)
}

func (u *userUsecase) Update(user domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) Delete(email string) error {
	return u.userRepo.Delete(email)
}
