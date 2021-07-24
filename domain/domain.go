package domain

type User struct {
	Email    string `json:"email" gorm:"primaryKey"`
	Password []byte `json:"password"` // Off course, HASHED
}

type UserUseCase interface {
	Create(user User) error
	Delete(email string) error
	Read(email string) (User, error)
}

type UserRepository interface {
	Create(user User) error
	Delete(email string) error
	Read(email string) (User, error)
}

type Session struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
}

type SessionUseCase interface {
	Create(session Session) error
	Delete(id string) error
	Read(id string) (Session, error)
}

type SessionRepository interface {
	Create(session Session) error
	Delete(id string) error
	Read(id string) (Session, error)
}
