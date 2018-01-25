package accountsbundle

import (
	"time"
	"log"
)

type Account struct {
	ID			int			`json:"id" db:"id"`
	Name		string		`json:"name" db:"name"`
	UserName	string		`json:"username" db:"username"`
	Balance		int64		`json:"balance" db:"balance"`
	Disabled	bool		`json:"disabled" db:"disabled"`
	Created		string		`json:"-" db:"created"`
}

func AddAccount(name, username string, balance int64) *Account {
	if balance < 0 { return nil }

	account := &Account{
		Name: name,
		UserName: username,
		Balance: balance,
		Disabled: false,
		Created: time.Now().Format(time.RFC3339),
	}

	if err := addAccountToDB(account); err != nil {
		log.Println(err)
		return nil
	}

	return account
}

func GetAccounts() *[]Account {
	accounts, err := getAccountsFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return accounts
}

func GetAccountByID(id int) *Account {
	account, err := getAccountByIDFromDB(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return account
}

func GetAccountByUsername(username string) *Account {
	account, err := getAccountByUsernameFromDB(username)
	if err != nil {
		log.Println(err)
		return nil
	}
	return account
}

func (a *Account) save() bool {
	if err := updateAccountInDB(*a); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (a *Account) IncreaseBalance(sum int64) bool {
	a.Balance += sum
	return a.save()
}

func (a *Account) DecreaseBalance(sum int64) bool {
	if sum > a.Balance {
		return false
	}
	a.Balance -= sum
	return a.save()
}

func (a *Account) Delete() bool {
	if err := deleteAccountFromDB(a.ID); err != nil {
		log.Println(err)
		return false
	}
	return true
}