package usecase

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/QQsha/sense-wallet/models"
	"github.com/QQsha/sense-wallet/usecase/mock"
)

// Get Balance

// test success
func TestGetBalance(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{UserID: userID, AmountUSD: 100}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	balance, err := wallet.GetBalance("1234")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if balance.AmountUSD != 100 || balance.UserID != "1234" {
		t.Errorf("unexpected balance: %v", balance)
	}
}

// user doesn't exist
func TestGetBalanceNoUser(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{}, sql.ErrNoRows
	}
	wallet := NewWallet(mockStore, log.Default())

	_, err := wallet.GetBalance("1234")
	if !errors.Is(err, userNotExistError) {
		t.Errorf("expected error: %v, got: %v", userNotExistError, err)
	}
}

// TRANSACTION

// test future time
func TestTransactionFutureTime(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	futureTime := time.Now().Add(time.Hour)
	transaction := models.Transaction{Time: futureTime}
	_, err := wallet.CreateTransaction(transaction)
	if !errors.Is(err, futureDateError) {
		t.Errorf("expected error: %v, got: %v", futureDateError, err)
	}
}

// test low balance euro
func TestTransactionLowBalance(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{AmountEuro: 10}, nil
	}
	wallet := NewWallet(mockStore, log.Default())
	currencies := []string{"EUR", "USD"}
	for _, curr := range currencies {
		transaction := models.Transaction{Currency: curr, Amount: 10.01, Type: "withdrawal"}
		_, err := wallet.CreateTransaction(transaction)
		if !errors.Is(err, lowBalanceError) {
			t.Errorf("expected error: %v, got: %v", lowBalanceError, err)
		}
	}
}

// success change +balance usd
func TestTransactionDepositUSD(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{AmountUSD: 100}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	transaction := models.Transaction{Currency: "USD", Amount: 10.01, Type: "deposit"}
	balance, err := wallet.CreateTransaction(transaction)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if balance.AmountUSD != 110.01 {
		t.Errorf("unexpected balance: %v", balance.AmountUSD)
	}
}

// success change +balance euro
func TestTransactionDepositEuro(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{AmountEuro: 100}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	transaction := models.Transaction{Currency: "EUR", Amount: 10.001, Type: "deposit"}
	balance, err := wallet.CreateTransaction(transaction)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if balance.AmountEuro != 110.001 {
		t.Errorf("unexpected balance: %v", balance.AmountEuro)
	}
}

// success change -balance usd
func TestTransactionWithdrawalUSD(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{AmountUSD: 50}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	transaction := models.Transaction{Currency: "USD", Amount: 10.001, Type: "withdrawal"}
	balance, err := wallet.CreateTransaction(transaction)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if balance.AmountUSD != 39.999 {
		t.Errorf("unexpected balance: %v", balance.AmountUSD)
	}

}

// success change -balance euro
func TestTransactionWithdrawalEuro(t *testing.T) {
	mockStore := mock.NewRepositoryMock()

	mockStore.GetBalanceFunc = func(userID string) (models.Balance, error) {
		return models.Balance{AmountEuro: 100}, nil
	}
	wallet := NewWallet(mockStore, log.Default())

	transaction := models.Transaction{Currency: "EUR", Amount: 10.01, Type: "withdrawal"}
	balance, err := wallet.CreateTransaction(transaction)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if balance.AmountEuro != 89.99 {
		t.Errorf("unexpected balance: %v", balance.AmountEuro)
	}

}
