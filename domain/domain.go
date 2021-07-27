package domain

type User struct {
	Email    string `json:"email" gorm:"primaryKey"`
	Password []byte `json:"password"` // Off course, HASHED
	LetterIDs []string `json:"letter_ids"`
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

type Letter struct {
	ID    string `json:"id" gorm:"primaryKey"`
	AuthorID string `json:"author_id"`
	Content string `json:"content"`
	IsSent bool `json:"is_sent"`
}

type LetterUseCase interface {
	Create(letter Letter) error
	Read(id string) (Letter, error)
	Update(letter Letter) error
	Delete(id string) error
	GetFirstUnsendLetter(AuthorID string) (Letter, error)
}

type LetterRepository interface {
	Create(letter Letter) error
	Read(id string) (Letter, error)
	Update(letter Letter) error
	Delete(id string) error
	GetFirstUnsendLetter(AuthorID string) (Letter, error)
}