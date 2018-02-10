package transactionsbundle

import (
	"log"

	"github.com/osoderholm/eletab-lite/app/bundles/accountsbundle"
	"github.com/osoderholm/eletab-lite/app/common"
)

// Empty helper struct for DBHelper interface
type helper struct{}

// Mandatory table creation function
func (helper helper) CreateTable(database *common.Database) error {
	qry := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			type INTEGER NOT NULL,
			account_id INTEGER NOT NULL,
			card_id TEXT NOT NULL,
			sum INT64 NOT NULL,
			accepted INTEGER NOT NULL,
			time DATETIME NOT NULL
			);`

	_, err := database.Exec(qry)

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

// Adds transaction to DB and gives transaction its ID
func addTransactionToDB(transaction *Transaction) error {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO transactions (
			type, account_id, card_id, sum, accepted, time
		) VALUES (
			:type, :account_id, :card_id, :sum, :accepted, :time
		);
		`
	res, err := db.NamedExec(query, transaction)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	transaction.ID = int(id)

	return nil
}

// Updates transaction details
func updateTransactionInDB(transaction Transaction) error {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	query := `
			UPDATE
				transactions
			SET
				type = :type,
				account_id = :account_id,
				card_id = :card_id,
				sum = :sum,
				accepted = :accepted
			WHERE transactions.id = :id;`

	_, err = db.NamedExec(query, &transaction)
	return err
}

// Gets from 'start' (including) to 'end' (excluding)
func getTransactionsByTimeFromDB(start, end string) (*[]Transaction, error) {
	var transactions []Transaction

	db, err := openDB()
	defer db.Close()
	if err != nil {
		return &transactions, err
	}

	query := `SELECT * FROM transactions 
				WHERE transactions.time >= ? AND transactions.time < ?`

	return &transactions, db.Select(&transactions, query, start, end)
}

// Gets by account, from 'start' (including) to 'end' (excluding)
func getTransactionsByAccountByTimeFromDB(account accountsbundle.Account, start, end string) (*[]Transaction, error) {
	var transactions []Transaction

	db, err := openDB()
	defer db.Close()
	if err != nil {
		return &transactions, err
	}

	query := `SELECT * FROM transactions 
				WHERE transactions.time >= ? AND transactions.time < ? AND account_id = ?`

	return &transactions, db.Select(&transactions, query, start, end, account.ID)
}
