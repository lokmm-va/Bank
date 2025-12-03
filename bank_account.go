package main

import (
	"fmt"
	"strings"
	"time"
)

type BankAccount struct {
	account *Account
	storage Storage
}

func NewBankAccount(account *Account, storage Storage) *BankAccount {
	return &BankAccount{
		account: account,
		storage: storage,
	}
}

func (ba *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	ba.account.Balance += amount
	transaction := Transaction{
		ID:          generateID(),
		Type:        Deposit,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Пополнение: %.2f", amount),
	}
	ba.account.Transactions = append(ba.account.Transactions, transaction)

	return ba.storage.SaveAccount(ba.account)
}

func (ba *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if ba.account.Balance < amount {
		return ErrInsufficientFunds
	}

	ba.account.Balance -= amount
	transaction := Transaction{
		ID:          generateID(),
		Type:        Withdraw,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Снятие: %.2f", amount),
	}
	ba.account.Transactions = append(ba.account.Transactions, transaction)

	return ba.storage.SaveAccount(ba.account)
}

func (ba *BankAccount) Transfer(to *Account, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if ba.account.Balance < amount {
		return ErrInsufficientFunds
	}

	if ba.account.ID == to.ID {
		return ErrSameAccountTransfer
	}

	ba.account.Balance -= amount
	transactionFrom := Transaction{
		ID:          generateID(),
		Type:        Transfer,
		Amount:      amount,
		Timestamp:   time.Now(),
		To:          to.ID,
		Description: fmt.Sprintf("Перевод на %s: %.2f", to.ID, amount),
	}
	ba.account.Transactions = append(ba.account.Transactions, transactionFrom)

	to.Balance += amount
	transactionTo := Transaction{
		ID:          generateID(),
		Type:        Receive,
		Amount:      amount,
		Timestamp:   time.Now(),
		From:        ba.account.ID,
		Description: fmt.Sprintf("Поступление от %s: %.2f", ba.account.ID, amount),
	}
	to.Transactions = append(to.Transactions, transactionTo)

	if err := ba.storage.SaveAccount(ba.account); err != nil {
		return err
	}

	return ba.storage.SaveAccount(to)
}

func (ba *BankAccount) GetBalance() float64 {
	return ba.account.Balance
}

func (ba *BankAccount) GetStatement() string {
	var sb strings.Builder

	sb.WriteString("================================\n")
	sb.WriteString(fmt.Sprintf("Счет: %s\n", ba.account.ID))
	sb.WriteString(fmt.Sprintf("Владелец: %s\n", ba.account.OwnerName))
	sb.WriteString(fmt.Sprintf("Баланс: %.2f\n", ba.account.Balance))
	sb.WriteString("Транзакции:\n")

	if len(ba.account.Transactions) == 0 {
		sb.WriteString("Нет транзакций\n")
	} else {
		for _, t := range ba.account.Transactions {
			sb.WriteString(fmt.Sprintf("%s %s\n",
				t.Timestamp.Format("02.01.2006 15:04"),
				t.Description))
		}
	}
	sb.WriteString("================================\n")

	return sb.String()
}
