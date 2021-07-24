package repository

import (
	"errors"
	"fmt"
	"sync"
	"user-auth/domain"
)

type userRepository struct {
	m sync.Map
}

func NewSyncMapUserRepository() domain.UserRepository {
	return &userRepository{}
}

func (u *userRepository) Create(user domain.User) error {
	u.m.Store(user.Email, user)
	value, ok := u.m.Load(user.Email)
	fmt.Println(value)
	fmt.Println(ok)
	return nil
}

func (u *userRepository) Delete(email string) error {
	u.m.Delete(email)
	return nil
}

func (u *userRepository) Check(email string) (domain.User, error) {
	value, ok := u.m.Load(email)
	fmt.Println(value)
	fmt.Println(ok)
	if ok {
		return value.(domain.User), nil
	}
	return domain.User{}, errors.New("user not found")
}
