package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*int, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp,
		encrypted_password varchar(150)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) (*int, error) {
	query := `insert into account
	(first_name, last_name, number, balance, created_at, encrypted_password)
	values ($1, $2, $3, $4, $5, $6) returning id`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var accountId int
	if err := stmt.QueryRow(a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt, a.EncryptedPassword).Scan(&accountId); err != nil {
		return nil, err
	}

	return &accountId, nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := "delete from account where id = $1"

	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(n int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", n)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", n)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt, &account.EncryptedPassword)
	if err != nil {
		return nil, err
	}
	return account, err
}
