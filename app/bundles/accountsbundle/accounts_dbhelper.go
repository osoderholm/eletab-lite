package accountsbundle

import (
	"log"
	"github.com/osoderholm/eletab-lite/app/common"
)

// Empty helper struct for DBHelper interface
type helper struct {}

// Mandatory table creation function
func (helper helper) CreateTable(database *common.Database) error {
	qryAccounts := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			balance INT64 CHECK(balance >= 0),
			level INTEGER NOT NULL,
			disabled INTEGER NOT NULL,
			created DATETIME NOT NULL
			);`

	_, err := database.Exec(qryAccounts)

	if err != nil {
		log.Println(err)
		return err
	}

	qryCards := `
		CREATE TABLE IF NOT EXISTS cards (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			card_id TEXT UNIQUE NOT NULL,
			account_id INTEGER NOT NULL,
			disabled INTEGER NOT NULL,
			added DATETIME NOT NULL
			);`

	_, err = database.Exec(qryCards)

	if err != nil {
		log.Println(err)
	}

	return err
}

// Wrapper function for opening common.Database
func openDB() (*common.Database, error) {
	db, err := common.OpenDB(helper{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// ********************************************************

// START Account DB

func addAccountToDB(account *Account) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `
		INSERT INTO accounts (
			name, username, password, balance, level, disabled, created
		) VALUES (
			:name, :username, :password, :balance, :level, :disabled, :created
		);
		`
	res, err := db.NamedExec(query, account)

	if err != nil { return err }

	id, err := res.LastInsertId()

	if err != nil { return err }

	account.ID = int(id)

	return nil
}

func getAccountsFromDB() (*[]Account, error) {
	var accounts []Account

	db, err := openDB()
	defer db.Close()
	if err != nil { return &accounts, err }

	query := `SELECT * FROM accounts;`

	return &accounts, db.Select(&accounts, query)
}

func getAccountByIDFromDB(accountID int) (*Account, error) {
	var account Account

	db, err := openDB()
	defer db.Close()
	if err != nil { return &account, err }

	query := `SELECT * FROM accounts WHERE accounts.id = ?;`

	return &account, db.Get(&account, query, accountID)
}

func getAccountByUsernameFromDB(username string) (*Account, error) {
	var account Account

	db, err := openDB()
	defer db.Close()
	if err != nil { return &account, err }

	query := `SELECT * FROM accounts WHERE accounts.username = ?;`

	return &account, db.Get(&account, query, username)
}

func updateAccountInDB(account Account) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `
			UPDATE
				accounts
			SET
				name = :name,
				username = :username,
				password = :password,
				balance = :balance,
				level = :level,
				disabled = :disabled
			WHERE accounts.id = :id;`

	_, err = db.NamedExec(query, &account)
	return err
}

func deleteAccountFromDB(accountID int) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `DELETE FROM accounts WHERE accounts.id = ?;`

	_, err = db.Exec(query, accountID)
	return err
}

// END Account DB

// ********************************************************

// START Card DB

func addCardToDB(card *Card) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `
		INSERT INTO cards (
			card_id, account_id, disabled, added
		) VALUES (
			:card_id, :account.id, :disabled, :added
		);
		`
	res, err := db.NamedExec(query, card)

	if err != nil { return err }

	id, err := res.LastInsertId()

	if err != nil { return err }

	card.ID = int(id)

	return nil
}

func getCardsFromDB() (*[]Card, error) {
	var cards []Card

	db, err := openDB()
	defer db.Close()
	if err != nil { return &cards, err }

	query := `SELECT
				cards.*,
				accounts.id "account.id",
				accounts.name "account.name",
				accounts.username "account.username",
				accounts.balance "account.balance",
				accounts.level "account.level",
				accounts.disabled "account.disabled",
				accounts.created "account.created"
			FROM
				cards JOIN accounts ON cards.account_id = accounts.id;`

	return &cards, db.Select(&cards, query)
}

func getCardsByAccountFromDB(account *Account) (*[]Card, error) {
	var cards []Card

	db, err := openDB()
	defer db.Close()
	if err != nil { return &cards, err }

	query := `SELECT
				cards.*,
				accounts.id "account.id",
				accounts.name "account.name",
				accounts.username "account.username",
				accounts.balance "account.balance",
				accounts.level "account.level",
				accounts.disabled "account.disabled",
				accounts.created "account.created"
			FROM
				cards JOIN accounts ON cards.account_id = accounts.id
			WHERE accounts.id = ?;`

	return &cards, db.Select(&cards, query, account.ID)
}

func getCardByCardIDFromDB(cardID string) (*Card, error) {
	var card Card

	db, err := openDB()
	defer db.Close()
	if err != nil { return &card, err }

	query := `SELECT
				cards.*,
				accounts.id "account.id",
				accounts.name "account.name",
				accounts.username "account.username",
				accounts.balance "account.balance",
				accounts.level "account.level",
				accounts.disabled "account.disabled",
				accounts.created "account.created"
			FROM
				cards JOIN accounts ON cards.account_id = accounts.id
			WHERE cards.card_id = ?;`

	return &card, db.Get(&card, query, cardID)
}

func deleteCardFromDB(cardID int) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `DELETE FROM cards WHERE cards.id = ?;`

	_, err = db.Exec(query, cardID)
	return err
}

// END Card DB

