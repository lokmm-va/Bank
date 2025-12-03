package main

import (
	"fmt"
	"time"
)

var transactionCounter = 0

func generateID() string {
	transactionCounter++
	return fmt.Sprintf("TXN%08d", transactionCounter)
}

func generateAccountID() string {
	return fmt.Sprintf("ACC%08d", time.Now().UnixNano()%100000000)
}
