package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto    *paseto.V2
	secretKey []byte
}

func NewPasetoMaker(secretKey string) (Maker, error) {
	if len(secretKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("The secret key length must be exactly %d characters long but giver %d", chacha20poly1305.KeySize, len(secretKey))
	}

	return &PasetoMaker{
		paseto:    paseto.NewV2(),
		secretKey: []byte(secretKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string,*Payload ,error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token, err :=maker.paseto.Encrypt(maker.secretKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.secretKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("token is invalid.")
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
