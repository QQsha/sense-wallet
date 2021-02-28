package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/QQsha/sense-wallet/models"
	"github.com/QQsha/sense-wallet/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type WalletHandler struct {
	use *usecase.Wallet
}

type response struct {
	Error struct {
		Code   int    `json:"code"`
		Detail string `json:"detail"`
	} `json:"error"`
}

func NewWalletHanlder(use *usecase.Wallet) *WalletHandler {
	return &WalletHandler{use: use}
}

func (u WalletHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fillResponse(400, err))
		return
	}

	// check if time have right format
	transaction.Time, err = time.Parse("2-Jan-06 15:04:05", transaction.TimePlaced)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fillResponse(400, err))
		return
	}
	// input validate
	validate := validator.New()
	err = validate.Struct(transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fillResponse(400, err))
		return
	}
	// make transaction
	balance, err := u.use.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fillResponse(500, err))
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (u WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	balance, err := u.use.GetBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(fillResponse(404, err))
		return
	}
	json.NewEncoder(w).Encode(balance)
}

func fillResponse(code int, err error) response {
	return response{
		Error: struct {
			Code   int    "json:\"code\""
			Detail string "json:\"detail\""
		}{Code: code, Detail: err.Error()},
	}
}
