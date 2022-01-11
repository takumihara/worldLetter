package domain

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
