package usecase

import (
	"github.com/tacomea/worldLetter/domain"
	"github.com/tacomea/worldLetter/domain/repository"
)

type UserUseCase interface {
	Create(user domain.User) error
	Read(email string) (domain.User, error)
	Update(user domain.User) error
	Delete(email string) error
}


type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(ur repository.UserRepository) UserUseCase {
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
