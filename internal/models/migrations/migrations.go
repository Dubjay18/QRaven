package migrations

import "qraven/internal/models"



func AuthMigrationModels() []interface{} {
	return []interface{}{
		&models.User{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}