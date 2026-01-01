package models

type Flight struct {
    ID        uint    `json:"id" gorm:"primaryKey"`
    Airline   string  `json:"airline"`
    From      string  `json:"from"`
    To        string  `json:"to"`
    Departure string  `json:"departure"`
    Arrival   string  `json:"arrival"`
    Price     float64 `json:"price"`
    City      string  `json:"city"`
    SeatsAvailable int `json:"seats_available"`
    Reservations []Reservation `json:"reservations" gorm:"foreignKey:FlightID"`
}
