package usecase

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/QQsha/sense-wallet/decimals"
	"github.com/QQsha/sense-wallet/models"
	"github.com/QQsha/sense-wallet/repository"
)

var (
	userNotExistError = errors.New("user doesn't exist")
	lowBalanceError   = errors.New("low balance")
	futureDateError   = errors.New("transcation date should be in the past")
)

type Wallet struct {
	store repository.StoreRepository
	log   *log.Logger
}

func NewWallet(store repository.StoreRepository, log *log.Logger) *Wallet {
	err := store.CreateTableBalance()
	if err != nil {
		log.Fatal(err)
	}
	return &Wallet{store: store, log: log}
}

func (wl Wallet) CreateTransaction(transaction models.Transaction) (models.Balance, error) {
	// transtaction time check
	if time.Until(transaction.Time) > 0 {
		wl.log.Println("wrong time in transaction")
		return models.Balance{}, futureDateError
	}

	balance, err := wl.store.GetBalance(transaction.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		wl.log.Println("cant get user balance")
		return balance, err
	}
	// amount switch to negative
	if transaction.Type == "withdrawal" {
		transaction.Amount = decimals.NegativeAmount(transaction.Amount)
	}
	// right balance choice
	switch transaction.Currency {
	case "EUR":
		newBalance := decimals.NewBalance(balance.AmountEuro, transaction.Amount)
		if newBalance < 0 {
			wl.log.Println("low euro balance")
			return balance, lowBalanceError
		}
		if err := wl.store.ChangeBalanceEuro(transaction.UserID, newBalance); err != nil {
			return balance, err
		}
		balance.AmountEuro = newBalance
	case "USD":
		newBalance := decimals.NewBalance(balance.AmountUSD, transaction.Amount)
		if newBalance < 0 {
			wl.log.Println("low usd balance")
			return balance, lowBalanceError
		}
		if err := wl.store.ChangeBalanceUSD(transaction.UserID, newBalance); err != nil {
			return balance, err
		}
		balance.AmountUSD = newBalance
	}

	return balance, nil
}

func (wl Wallet) GetBalance(userID string) (models.Balance, error) {
	balance, err := wl.store.GetBalance(userID)
	if errors.Is(err, sql.ErrNoRows) {
		wl.log.Println("cant get user balance")
		return balance, userNotExistError
	}
	return balance, err
}
