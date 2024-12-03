package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	notificationService "qraven/services/notification"

	"log"
	"qraven/internal/config"
	"qraven/internal/models/migrations"
	"qraven/pkg/repository/storage"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/pkg/repository/storage/redis"
	"qraven/pkg/router"
	"qraven/utils"
)

// @title QRaven API
// @version 1.0
// @description This is a sample QRaven API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8019
// @BasePath /api/v1
func main() {
	logger := utils.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app
	configuration := config.Setup(logger, "./app")
	postgresql.ConnectToDatabase(logger, configuration.Database)
	redis.ConnectToRedis(logger, configuration.Redis)
	validatorRef := validator.New()

	db := storage.Connection()

	if configuration.Database.Migrate {
		migrations.RunAllMigrations(db)
		// call the seed function
		// seed.SeedDatabase(db.Postgresql)
	}
	c := cron.New()
	c.AddFunc("@daily", func() {
		err := notificationService.CleanupExpiredTokens(db)
		if err != nil {
			log.Println("Error cleaning up expired tokens:", err)
		} else {
			log.Println("Expired tokens cleaned up successfully")
		}
	})
	c.Start()

	select {} // Keep the application running

	r := router.Setup(logger, validatorRef, db, &configuration.App)
	utils.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	log.Fatal(r.Run(":" + configuration.Server.Port))
}
