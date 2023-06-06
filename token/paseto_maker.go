package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// check the length of the symmetric key
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key length: %d", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}