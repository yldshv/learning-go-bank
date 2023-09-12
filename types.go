package main

import "math/rand"

type Account struct {
	ID        int     `json:"id"`
	FirstName string  `json:"fistName"`
	LastName  string  `json:"lastName"`
	Number    int64   `json:"number"`
	Balance   float64 `json:"balance"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int63n(100000000000),
	}
}
