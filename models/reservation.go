package models

type Reservation struct {
	ID uint `json:"id" gorm:"primaryKey"`

	UserID  string `json:"userId"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
	Status  string `json:"status"`
	HotelID string `json:"hotel_id"`
	Hotel   Hotel  `json:"hotel" gorm:"foreignKey:HotelID"`

	FlightID *string `json:"flight_id"` // OPTIONAL
	Flight   *Flight `json:"flight" gorm:"foreignKey:FlightID"`

	CheckIn    string  `json:"check_in"`
	CheckOut   string  `json:"check_out"`
	TotalPrice float64 `json:"total_price"`
}
