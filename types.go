package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TransferRequest struct {
	ToAccount int     `json:"toAccount"`
	Amount    float64 `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
	EncryptedPassword string    `json:"-"`
}

func NewAccount(firstName string, lastName string, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            rand.Int63n(1000000),
		CreatedAt:         time.Now().UTC(),
		EncryptedPassword: string(encpw),
	}, nil
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

func (a *Account) PWisValid(p string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(p)); err != nil {
		return false
	}
	return true
}
