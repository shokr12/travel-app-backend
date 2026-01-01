package models

type VisaApplication struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	UserID         uint
	User           User `gorm:"foreignKey:UserID"`
	VisaType       string
	Destination    string
	TravelDate     string
	PassportNumber string
	Nationality    string
	Country        string
	PassportURL    string
	Status         string // pending, approved, rejected
}
type CreateVisaRequest struct {
	UserID         uint   `json:"userId" binding:"required"`
	VisaType       string `json:"visa_type" binding:"required,min=2,max=50"`
	Destination    string `json:"destination" binding:"required,min=2,max=100"`
	TravelDate     string `json:"travel_date" binding:"required"`
	PassportNumber string `json:"passport_number" binding:"required,min=6,max=20"`
	Nationality    string `json:"nationality" binding:"required,min=2,max=50"`
}

type UpdateVisaRequest struct {
	VisaType       string `json:"visa_type" binding:"omitempty,min=2,max=50"`
	Destination    string `json:"destination" binding:"omitempty,min=2,max=100"`
	TravelDate     string `json:"travel_date" binding:"omitempty"`
	PassportNumber string `json:"passport_number" binding:"omitempty,min=6,max=20"`
	Nationality    string `json:"nationality" binding:"omitempty,min=2,max=50"`
}
