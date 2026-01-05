package models

type SupportTicket struct {
	ID uint `json:"id" gorm:"primaryKey"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	Subject string `json:"subject" binding:"required,min=10,max=2000"`
	Message string `json:"message" binding:"required,min=10,max=2000"`
	Status  string `json:"status" binding:"required,oneof=open in_progress resolved closed"`
}
