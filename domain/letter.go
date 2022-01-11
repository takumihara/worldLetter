package domain

type Letter struct {
	ID    string `json:"id" gorm:"primaryKey"`
	AuthorID string `json:"author_id"`
	ReceiverID string `json:"receiver_id"`
	Content string `json:"content"`
	IsSent bool `json:"is_sent"`
}

type LetterUseCase interface {
	Create(letter Letter) error
	Read(id string) (Letter, error)
	Update(letter Letter) error
	Delete(id string) error
	GetAll() ([]Letter, error)
	GetFirstUnsendLetter(authorID string) (Letter, error)
	GetLettersByAuthorID(authorID string) (string, error)
	GetLettersByReceiverID(receiverID string) (string, error)
}

type LetterRepository interface {
	Create(letter Letter) error
	Read(id string) (Letter, error)
	Update(letter Letter) error
	Delete(id string) error
	GetAll() ([]Letter, error)
	GetFirstUnsendLetter(authorID string) (Letter, error)
	GetLettersByAuthorID(authorID string) (string, error)
	GetLettersByReceiverID(receiverID string) (string, error)
}
