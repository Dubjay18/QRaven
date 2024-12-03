package notificationService

import (
	"errors"
	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"time"
)

func ExpoNotify(c *gin.Context, db *storage.Database) error {
	pushToken, err := expo.NewExponentPushToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]")
	if err != nil {
		return err
	}

	// Create a new Expo SDK client
	client := expo.NewPushClient(nil)

	// Publish message
	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     "This is a test notification",
			Data:     map[string]string{"withSome": "data"},
			Sound:    "default",
			Title:    "Notification Title",
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		return err
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		return errors.New("failed to send notification")
	}
	return nil
}

func SaveExpoToken(c *gin.Context, db storage.Database, token string) error {
	expoPushToken := models.ExpoPushToken{
		Token: token,
	}
	err := db.Postgresql.Create(&expoPushToken).Error
	if err != nil {
		return err
	}
	return nil
}

func CleanupExpiredTokens(db *storage.Database) error {
	return db.Postgresql.Where("expires_at < ?", time.Now()).Delete(&models.ExpoPushToken{}).Error
}
