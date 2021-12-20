package domain

type Letter struct {
	ID         string `json:"id" gorm:"primaryKey"`
	AuthorID   string `json:"author_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	IsSent     bool   `json:"is_sent"`
}
