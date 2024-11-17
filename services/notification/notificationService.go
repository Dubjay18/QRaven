package notification

import (
	"errors"
	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"qraven/pkg/repository/storage"
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
