package domain

type User struct {
	Email    string `json:"email" gorm:"primaryKey"`
	Password []byte `json:"password"` // Off course, HASHED
}
