package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// docker run --name some-postgres -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error) // Fixed: Should return a slice of *Account and an error
	GetAccountById(int) (*Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgressStore, error) {
	// connStr := "user=postgres dbname=gobank sslmode=disable"
	connStr := "user=postgres password=root dbname=gobank sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgressStore{
		db: db,
	}, nil
}

// createAccountTable creates the account table if it doesn't exist
func (s *PostgressStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number INT,
		balance INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgressStore) CreateAccount(acc *Account) error {
	query := `
INSERT INTO account (first_name, last_name, number, balance, created_at)
VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgressStore) UpdateAccount(acc *Account) error {
	query := `UPDATE account SET first_name=$1,last_name=$2,number=$3,balance=$4 where id=$5`
	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.ID)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgressStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}
func (s *PostgressStore) GetAccountById(id int) (*Account, error) {
	// Correct SQL query with proper column selection
	rows, err := s.db.Query("select * from account where id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanFromDB(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)

}

func (s *PostgressStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		account := new(Account)
		account, err := scanFromDB(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanFromDB(row *sql.Rows) (*Account, error) {
	account := new(Account)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err

}
