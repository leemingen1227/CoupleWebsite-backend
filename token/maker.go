package token

import (
	"time"

	"github.com/google/uuid"
)

//Maker is an interface for creating and validating tokens.
type Maker interface {
	//CreateToken creates a new token for the given user ID and duration.
	CreateToken(userID uuid.UUID, duration time.Duration) (string, *Payload, error)
	//VerifyToken verifies the given token and returns the payload.
	VerifyToken(token string) (*Payload, error)
}
