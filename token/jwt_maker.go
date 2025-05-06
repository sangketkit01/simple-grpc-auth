package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("Secret key must be at least 32 characters long but given %d", len(secretKey))
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error){
	payload, err := NewPayload(username, duration) 
	if err != nil{
		return "", err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwt.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload,error){
	keyFunc := func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok{
			return nil, fmt.Errorf("Invalid token.")
		}

		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil{
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, fmt.Errorf("token is invalid")){
			return nil, fmt.Errorf("token is invalid")
		}

		return nil, fmt.Errorf("token is invalid")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok{
		return nil, fmt.Errorf("token is invalid")
	}

	err = payload.Valid()
	if err != nil{
		return nil, err
	}

	return payload, nil
}