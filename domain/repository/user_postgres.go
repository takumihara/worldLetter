package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/gorm"
)

type userRepositoryPG struct {
	db *gorm.DB
}

func NewUserRepositoryPG(db *gorm.DB) UserRepository {
	return &userRepositoryPG{
		db: db,
	}
}

func (u *userRepositoryPG) Create(user domain.User) error {
	result := u.db.Create(user)
	if result.Error != nil {
		return errors.New("unexpected error while creating user")
	}
	return nil
}

func (u *userRepositoryPG) Read(email string) (domain.User, error) {
	var user domain.User
	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return domain.User{}, errors.New("unexpected error while retrieving user")
	}
	return user, nil
}

func (u *userRepositoryPG) Update(user domain.User) error {
	u.db.Save(&user)
	return nil
}

func (u *userRepositoryPG) Delete(email string) error {
	result := u.db.Delete(&domain.User{}, email)
	if result.Error != nil {
		return errors.New("unexpected error while deleting user")
	}
	return nil
}
