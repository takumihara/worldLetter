package domain

type Session struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
}
