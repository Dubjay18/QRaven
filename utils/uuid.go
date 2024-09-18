package utils

import (
	"strings"

	"github.com/gofrs/uuid"
)

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func GenerateTicketId() string {
	// initials := GetInitialsFromEventName(eventName)
	id, _ := uuid.NewV7()
	return "qrv" + id.String()[:6]
}
func IsValidUUID(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}

func GetInitialsFromEventName(eventName string) string {
	words := strings.Fields(eventName)
	initials := ""
	for _, word := range words {
		initials += strings.ToUpper(string(word[0]))
	}
	return initials
}
