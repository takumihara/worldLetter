package domain

type User struct {
	Email    string `json:"email" gorm:"primaryKey"`
	Password []byte `json:"password"` // Off course, HASHED
}

type UserUseCase interface {
	Create(user User) error
	Read(email string) (User, error)
	Update(user User) error
	Delete(email string) error
}

type UserRepository interface {
	Create(user User) error
	Read(email string) (User, error)
	Update(user User) error
	Delete(email string) error
}
