package utils

import "github.com/gofrs/uuid"

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func GenerateTicketId(initials string) string {
	id, _ := uuid.NewV7()
	return initials + id.String()[:6]
}
func IsValidUUID(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}