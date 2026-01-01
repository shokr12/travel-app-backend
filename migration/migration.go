package migration

import (
	"Visa/config"
	"Visa/models"
)

func Migrate() {
	db := config.ConnectToDB()

	err := db.AutoMigrate(
		&models.User{},
		&models.Flight{},
		&models.Reservation{},
		&models.Hotel{},
		&models.VisaApplication{},
		&models.SupportTicket{},
	)

	if err != nil {
		panic("failed to migrate database")
	}
}
