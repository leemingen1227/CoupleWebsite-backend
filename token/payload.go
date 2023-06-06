package token

import (
	"github.com/google/uuid"
	"time"
	"errors"
)


var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

//Payload is the data encoded in a token.
type Payload struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

//Create new payload
func NewPayload(userID uuid.UUID, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	payload := &Payload{
		ID: id,
		UserID: userID,
		IssuedAt: now,
		ExpiresAt: now.Add(duration),
	}
	return payload, nil
}

//The function Valid check if the token is expired
func (p *Payload) Valid() error {
	if p.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
