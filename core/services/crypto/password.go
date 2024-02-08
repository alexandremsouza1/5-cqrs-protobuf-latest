package crypto

import (
	"context"

	"main.go/services/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx context.Context, pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, err)
	}
	return string(hash)
}

func CheckPassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}
