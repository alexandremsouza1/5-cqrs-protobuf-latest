package crypto

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IssuedAt   time.Time `json:"iat"`
	ExpiredAt  time.Time `json:"exp"`
	Iss        string    `json:"iss"`
	Picture    string    `json:"picture"`
	Authorized bool      `json:"authorized"`
	Email      string    `json:"email"`
	Roles      []string  `json:"roles"`
	Parties    []string  `json:"parties"`
	Companies  []string  `json:"companies"`
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func CreateToken(userName string, userEmail string, userPicture string) (string, error) {
	secret := os.Getenv("CORE_JWT_ACCESS_SECRET")
	duration := time.Hour * 8
	payload, err := NewPayload(userName, userEmail, userPicture, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(secret))
}

// get jwt token from metadata
func GetJWT(ctx context.Context) (string, error) {
	token, err := GetMetadata(ctx, "authorization")
	if err != nil {
		return "", err
	}
	// Checking if access token is missing
	if strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("Token is empty")
	}
	return token, nil
}

func GetToken(ctx context.Context) (*Payload, error) {
	token, err := GetJWT(ctx)
	if err != nil {
		return nil, err
	}
	secret := os.Getenv("CORE_JWT_ACCESS_SECRET")
	jwtString := token
	if strings.Contains(jwtString, "Bearer") {
		jwt := strings.Split(token, "Bearer ")
		if len(token) < 2 {
			return nil, ErrInvalidToken
		}
		jwtString = jwt[1]
	}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}
	jwtToken, err := jwt.ParseWithClaims(jwtString, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func VerifyToken(token string) (*Payload, error) {
	secret := os.Getenv("CORE_JWT_ACCESS_SECRET")
	jwtString := token
	if strings.Contains(jwtString, "Bearer") {
		jwt := strings.Split(token, "Bearer ")
		if len(token) < 2 {
			return nil, ErrInvalidToken
		}
		jwtString = jwt[1]
	}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}
	jwtToken, err := jwt.ParseWithClaims(jwtString, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewPayload(userName string, userEmail string, userPicture string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:         tokenID,
		Name:       userName,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().Add(duration),
		Iss:        "investcore.io",
		Picture:    userPicture,
		Authorized: true,
		Email:      userEmail,
	}
	return payload, nil
}

func NewJWTRandomKey() error {
	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		// really, what are you gonna do if randomness failed?
		return err
	}
	base64Key := base64.StdEncoding.EncodeToString([]byte(key))
	fmt.Printf("JWT Secret:       %s\n", base64Key)
	return nil
}

func GetMetadata(ctx context.Context, name string) (string, error) {
	headers, _ := metadata.FromIncomingContext(ctx)
	header := strings.Join(headers[name], " ")

	// Checking if access header is missing
	if strings.TrimSpace(header) == "" {
		return "", fmt.Errorf("Header key not found {%v}", name)
	}

	return header, nil
}
