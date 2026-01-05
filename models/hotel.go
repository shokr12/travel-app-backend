package models

type Hotel struct {
	ID               uint `json:"id" gorm:"primaryKey"`
	Name             string
	Location         string
	Description      string
	PricePerNight    float64
	SeatsAvailable   int
	City             string
	Address          string
	CheckInDate      string
	CheckOutDate     string
	FreeCancellation bool
	AvailableRooms   int

	Status       string
	Reservations []Reservation `gorm:"foreignKey:HotelID"`
	UserID uint `gorm:"column:user_id" json:"user_id"`
}
