package accountsbundle

import (
	"time"
	"log"
	"github.com/osoderholm/eletab-lite/app/crypt"
)

// User account. Contains balance and general account stuff.
// Cards are added to accounts, aka, these.
type Account struct {
	ID			int				`json:"id" db:"id"`
	Name		string			`json:"name" db:"name"`
	UserName	string			`json:"username" db:"username"`
	Password	string			`json:"-" db:"password"`
	Balance		int64			`json:"balance" db:"balance"`
	Level 		AccountLevel	`json:"level" db:"level"`
	Disabled	bool			`json:"disabled" db:"disabled"`
	Created		string			`json:"-" db:"created"`
}

// AccountLevel defines the level of rights for account.
// If following constants are edited, make sure to have a great
// DB migration plan, as these are saved there.
type AccountLevel int
const (
	LevelDefault 	= 0
	LevelModerator 	= 1
	LevelAdmin		= 2
	LevelSuperAdmin = 3
)

// Create new account, save it to DB and return it.
// Returns nil if unsuccessful.
// Observe that the password passed to this function needs to be in plaintext,
// as it is hashed within the function!
func AddAccount(name, username, password string, balance int64, level AccountLevel) *Account {
	if balance < 0 { return nil }

	passHash, err := crypt.Encrypt(password)

	if err != nil {
		log.Println(err)
		return nil
	}

	account := &Account{
		Name: name,
		UserName: username,
		Password: passHash,
		Balance: balance,
		Level: level,
		Disabled: false,
		Created: time.Now().Format(time.RFC3339),
	}

	if err := addAccountToDB(account); err != nil {
		log.Println(err)
		return nil
	}

	return account
}

// Returns all accounts. Use with care, although the passwords are hashed and not JSON marshaled.
func GetAccounts() *[]Account {
	accounts, err := getAccountsFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return accounts
}

// Get account by account ID. Nil if not found.
func GetAccountByID(id int) *Account {
	account, err := getAccountByIDFromDB(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return account
}

// Get account by its username. Nil if not found.
func GetAccountByUsername(username string) *Account {
	account, err := getAccountByUsernameFromDB(username)
	if err != nil {
		log.Println(err)
		return nil
	}
	return account
}

// Saves account data to DB.
func (a *Account) Save() bool {
	if err := updateAccountInDB(*a); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Increases account balance by 'sum' and saves to DB. Returns success.
func (a *Account) IncreaseBalance(sum int64) bool {
	a.Balance += sum
	return a.Save()
}

// Decreases account balance by 'sum' and saves to DB. Returns success.
func (a *Account) DecreaseBalance(sum int64) bool {
	if sum > a.Balance {
		return false
	}
	a.Balance -= sum
	return a.Save()
}

// Deletes the account from DB. Returns success.
// USE WITH CARE!
func (a *Account) Delete() bool {
	if err := deleteAccountFromDB(a.ID); err != nil {
		log.Println(err)
		return false
	}
	return true
}