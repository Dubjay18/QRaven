package postgresql

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"qraven/internal/config"
	"qraven/pkg/repository/storage"
	"qraven/utils"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

func ConnectToDatabase(logger *utils.Logger, configDatabases config.Database) *gorm.DB {
	dbsCV := configDatabases
	utils.LogAndPrint(logger, "connecting to database")
	connectedDB := connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.DB_NAME, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)

	utils.LogAndPrint(logger, "connected to database")

	storage.DB.Postgresql = connectedDB

	return connectedDB
}

func connectToDb(host, user, password, dbname, port, sslmode, timezone string, logger *utils.Logger) *gorm.DB {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utils.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utils.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
fmt.Println(port)
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", host, user, password, dbname, port, sslmode, timezone)
fmt.Println(dsn)
	newLogger := lg.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		lg.Config{
			LogLevel:                  lg.Error, // Log level
			IgnoreRecordNotFoundError: true,     // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		utils.LogAndPrint(logger, fmt.Sprintf("connection to %v db failed with: %v", dbname, err))
		panic(err)

	}

	utils.LogAndPrint(logger, fmt.Sprintf("connected to %v db", dbname))
	// db = db.Debug() //database debug mode
	return db
}