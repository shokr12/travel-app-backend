package models

type User struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Hide in JSON
	Role     string `json:"role"`

	Reservations     []Reservation     `json:"reservations" gorm:"foreignKey:UserID"`
	VisaApplications []VisaApplication `json:"visa_applications" gorm:"foreignKey:UserID"`
	SupportTickets   []SupportTicket   `json:"support_tickets" gorm:"foreignKey:UserID"`
}
