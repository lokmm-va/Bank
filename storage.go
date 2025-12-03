package main

type InMemoryStorage struct {
	accounts map[string]*Account
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		accounts: make(map[string]*Account),
	}
}

func (s *InMemoryStorage) SaveAccount(account *Account) error {
	s.accounts[account.ID] = account
	return nil
}

func (s *InMemoryStorage) LoadAccount(accountID string) (*Account, error) {
	account, exists := s.accounts[accountID]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *InMemoryStorage) GetAllAccounts() ([]*Account, error) {
	var accounts []*Account
	for _, account := range s.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}
