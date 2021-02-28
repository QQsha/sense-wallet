package mock

import "github.com/QQsha/sense-wallet/models"

type StoreMock struct {
	GetBalanceFunc         func(userID string) (models.Balance, error)
	ChangeBalanceEuroFunc  func(userID string, amount float64) error
	ChangeBalanceUSDFunc   func(userID string, amount float64) error
	CreateTableBalanceFunc func() error
}

func NewRepositoryMock() *StoreMock {
	return &StoreMock{}
}

func (m StoreMock) GetBalance(userID string) (models.Balance, error) {
	return m.GetBalanceFunc(userID)
}

func (m StoreMock) ChangeBalanceEuro(userID string, amount float64) error {
	return nil
}

func (m StoreMock) ChangeBalanceUSD(userID string, amount float64) error {
	return nil
}
func (m StoreMock) CreateTableBalance() error {
	return nil
}
