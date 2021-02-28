package main

import (
	"log"
	"net/http"

	"github.com/QQsha/sense-wallet/handlers"
	"github.com/QQsha/sense-wallet/repository"
	"github.com/QQsha/sense-wallet/usecase"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, teardown, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer teardown()

	store := repository.NewStore(db)
	logger := log.Default()

	wallet := usecase.NewWallet(store, logger)

	handler := handlers.NewWalletHanlder(wallet)

	r := mux.NewRouter()

	r.HandleFunc("/transaction", handler.CreateTransaction).Methods("POST")
	r.HandleFunc("/balance/{user_id}", handler.GetBalance).Methods("POST")
	logger.Println("Wallet started")
	log.Fatal(http.ListenAndServe(":8080", r))
}
