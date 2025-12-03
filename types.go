package main

import (
	"errors"
	"time"
)

var (
	ErrInsufficientFunds   = errors.New("недостаточно средств")
	ErrInvalidAmount       = errors.New("некорректная сумма")
	ErrAccountNotFound     = errors.New("счет не найден")
	ErrSameAccountTransfer = errors.New("перевод на тот же счет")
)

type TransactionType string

const (
	Deposit  TransactionType = "ПОПОЛНЕНИЕ"
	Withdraw TransactionType = "СНЯТИЕ"
	Transfer TransactionType = "ПЕРЕВОД"
	Receive  TransactionType = "ПОЛУЧЕНИЕ"
)

type Transaction struct {
	ID          string
	Type        TransactionType
	Amount      float64
	Timestamp   time.Time
	From        string
	To          string
	Description string
}

type Account struct {
	ID           string
	OwnerName    string
	Balance      float64
	Transactions []Transaction
	CreatedAt    time.Time
}

type AccountService interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	Transfer(to *Account, amount float64) error
	GetBalance() float64
	GetStatement() string
}

type Storage interface {
	SaveAccount(account *Account) error
	LoadAccount(accountID string) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}
