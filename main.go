package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func printMenu() {
	fmt.Println("\n БАНКОВСКИЕ ОПЕРАЦИИ ")
	fmt.Println("1. Создать счет")
	fmt.Println("2. Пополнить")
	fmt.Println("3. Снять")
	fmt.Println("4. Перевести")
	fmt.Println("5. Баланс")
	fmt.Println("6. Выписка")
	fmt.Println("7. Все счета")
	fmt.Println("8. Выход")
	fmt.Print("Выберите: ")
}

func inputFloat(prompt string) (float64, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, ErrInvalidAmount
	}
	return value, nil
}

func inputString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	storage := NewInMemoryStorage()
	var currentAccount *BankAccount
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\033[H\033[2J")
		printMenu()

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Println("\nСоздание счета")
			ownerName := inputString("Имя владельца: ")
			if ownerName == "" {
				fmt.Println("Ошибка: имя пустое")
				break
			}

			account := &Account{
				ID:        generateAccountID(),
				OwnerName: ownerName,
				Balance:   0,
				CreatedAt: time.Now(),
			}

			if err := storage.SaveAccount(account); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Printf("Счет создан: %s\n", account.ID)
				currentAccount = NewBankAccount(account, storage)
			}

			fmt.Print("Нажмите Enter чтобы вернуться в меню")
			reader.ReadString('\n')

		case "2":
			fmt.Println("\nПополнение")
			if currentAccount == nil {
				accountID := inputString("ID счета: ")
				account, err := storage.LoadAccount(accountID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					currentAccount = NewBankAccount(account, storage)
				}
			}

			if currentAccount != nil {
				amount, err := inputFloat("Сумма: ")
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					if err := currentAccount.Deposit(amount); err != nil {
						fmt.Printf("Ошибка: %v\n", err)
					} else {
						fmt.Printf("Успешно. Баланс: %.2f\n", currentAccount.GetBalance())
					}
				}
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "3":
			fmt.Println("\nСнятие")
			if currentAccount == nil {
				accountID := inputString("ID счета: ")
				account, err := storage.LoadAccount(accountID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					currentAccount = NewBankAccount(account, storage)
				}
			}

			if currentAccount != nil {
				amount, err := inputFloat("Сумма: ")
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					if err := currentAccount.Withdraw(amount); err != nil {
						fmt.Printf("Ошибка: %v\n", err)
					} else {
						fmt.Printf("Успешно. Баланс: %.2f\n", currentAccount.GetBalance())
					}
				}
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "4":
			fmt.Println("\nПеревод")
			if currentAccount == nil {
				accountID := inputString("ID вашего счета: ")
				account, err := storage.LoadAccount(accountID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
					break
				}
				currentAccount = NewBankAccount(account, storage)
			}

			if currentAccount != nil {
				targetID := inputString("ID получателя: ")
				targetAccount, err := storage.LoadAccount(targetID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
					break
				}

				amount, err := inputFloat("Сумма: ")
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					if err := currentAccount.Transfer(targetAccount, amount); err != nil {
						fmt.Printf("Ошибка: %v\n", err)
					} else {
						fmt.Printf("Перевод выполнен. Новый баланс: %.2f\n", currentAccount.GetBalance())
					}
				}
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "5":
			fmt.Println("\nБаланс")
			if currentAccount == nil {
				accountID := inputString("ID счета: ")
				account, err := storage.LoadAccount(accountID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					currentAccount = NewBankAccount(account, storage)
				}
			}

			if currentAccount != nil {
				fmt.Printf("Счет: %s\n", currentAccount.account.ID)
				fmt.Printf("Владелец: %s\n", currentAccount.account.OwnerName)
				fmt.Printf("Баланс: %.2f\n", currentAccount.GetBalance())
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "6":
			fmt.Println("\nВыписка")
			if currentAccount == nil {
				accountID := inputString("ID счета: ")
				account, err := storage.LoadAccount(accountID)
				if err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					currentAccount = NewBankAccount(account, storage)
				}
			}

			if currentAccount != nil {
				fmt.Println(currentAccount.GetStatement())
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "7":
			fmt.Println("\nВсе счета")
			accounts, err := storage.GetAllAccounts()
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else if len(accounts) == 0 {
				fmt.Println("Нет счетов")
			} else {
				for _, acc := range accounts {
					fmt.Printf("%s: %s - %.2f\n",
						acc.ID, acc.OwnerName, acc.Balance)
				}
			}

			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')

		case "8":
			fmt.Println("\nВыход...")
			return

		default:
			fmt.Println("Неверный выбор")
			fmt.Print("Нажмите Enter...")
			reader.ReadString('\n')
		}
	}
}
