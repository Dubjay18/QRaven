package migrations

import "qraven/internal/models"



func AuthMigrationModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Event{},
		&models.Ticket{},
		&models.Payments{},
		&models.Notification{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}