package repository

import "github.com/tacomea/worldLetter/domain"

type UserRepository interface {
	Create(user domain.User) error
	Read(email string) (domain.User, error)
	Update(user domain.User) error
	Delete(email string) error
}
