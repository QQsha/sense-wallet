package repository

import (
	"database/sql"

	"github.com/QQsha/sense-wallet/models"
)

type StoreRepository interface {
	GetBalance(userID string) (models.Balance, error)
	ChangeBalanceEuro(userID string, amount float64) error
	ChangeBalanceUSD(userID string, amount float64) error
	CreateTableBalance() error
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func ConnectDB() (*sql.DB, func() error, error) {
	db, err := sql.Open("sqlite3", "./wallet.db")
	if err != nil {
		return db, nil, err
	}
	return db, db.Close, nil
}

func (s Store) CreateTransaction(transaction models.Transaction) error {
	return nil
}

func (s Store) GetBalance(userID string) (models.Balance, error) {
	var balance models.Balance
	err := s.db.QueryRow(`SELECT * FROM balance WHERE user_id=?`, userID).Scan(&balance.UserID, &balance.AmountEuro, &balance.AmountUSD)

	return balance, err
}

func (s Store) ChangeBalanceEuro(userID string, amount float64) error {
	exec, err := s.db.Prepare(
		`INSERT INTO balance (user_id, amount_euro) 
		VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET amount_euro=?`)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Stmt(exec).Exec(userID, amount, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s Store) ChangeBalanceUSD(userID string, amount float64) error {
	exec, err := s.db.Prepare(
		`INSERT INTO balance (user_id, amount_usd) 
		VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET amount_usd=?`)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Stmt(exec).Exec(userID, amount, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s Store) CreateTableBalance() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS "balance" (
		"user_id" VARCHAR(64) PRIMARY KEY, 
		"amount_euro" REAL DEFAULT 0, 
		"amount_usd" REAL DEFAULT 0)`,
	)
	return err
}
