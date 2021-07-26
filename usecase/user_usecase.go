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
	err := u.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Read(email string) (domain.User, error) {
	user, err := u.userRepo.Read(email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userUsecase) Update(user domain.User) error {
	err := u.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Delete(email string) error {
	err := u.userRepo.Delete(email)
	if err != nil {
		return err
	}
	return nil
}
