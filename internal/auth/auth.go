// Package auth is used for all authentication
package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	p, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("error hashing passowrd: %v", err)
		return "", err
	}
	return p, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("error comparing password and hash: %v", err)
		return false, err
	}
	return match, nil
}
