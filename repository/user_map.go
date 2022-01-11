package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"sync"
)

type userRepository struct {
	m sync.Map
}

func NewSyncMapUserRepository() domain.UserRepository {
	return &userRepository{}
}

func (u *userRepository) Create(user domain.User) error {
	u.m.Store(user.Email, user)
	_, ok := u.m.Load(user.Email)
	if !ok {
		return errors.New("user not created")
	}
	return nil
}

func (u *userRepository) Read(email string) (domain.User, error) {
	value, ok := u.m.Load(email)
	if ok {
		return value.(domain.User), nil
	}
	return domain.User{}, errors.New("user not found")
}

func (u *userRepository) Update(user domain.User) error {
	u.m.Store(user.Email, user)
	_, ok := u.m.Load(user.Email)
	if !ok {
		return errors.New("user not created")
	}
	return nil
}

func (u *userRepository) Delete(email string) error {
	u.m.Delete(email)
	return nil
}
