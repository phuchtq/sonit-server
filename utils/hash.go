package utils

import (
	"errors"
	"log"
	"sonit_server/constant/noti"

	"golang.org/x/crypto/bcrypt"
)

func ToHashString(src string, logger *log.Logger) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(src), 10)

	if err != nil {
		logger.Println("Error while generating to hash string - " + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	return string(bytes), nil
}

func IsHashStringMatched(inputtedStr, hashedStr string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(inputtedStr)) == nil
}
