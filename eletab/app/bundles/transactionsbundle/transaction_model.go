package transactionsbundle

import (
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/accountsbundle"
	"time"
	"log"
)

// START TransactionManager

// Generic transaction manager for structure
type TransactionManager struct {}

// Creates and returns a TransactionManager
func NewManager() *TransactionManager {
	return &TransactionManager{}
}

// Removes from accounts balance in form of a purchase. Returns created transaction with details
func (tm *TransactionManager) MakePurchase(account *accountsbundle.Account, cardID string, sum int64) *Transaction {
	accepted := account.Balance >= sum && !account.Disabled

	transaction := makeTransaction(TypePurchase, account, cardID, sum, accepted)

	if transaction != nil {
		if !account.DecreaseBalance(sum) {
			transaction.Accepted = false
			transaction.save()
		}
	}

	return transaction
}

// Adds to accounts balance. Returns created transaction with details
func (tm *TransactionManager) MakeInsert(account *accountsbundle.Account, sum int64) *Transaction {
	transaction := makeTransaction(TypeInsert, account, "", sum, true)

	if transaction != nil {
		if !account.IncreaseBalance(sum) {
			transaction.Accepted = false
			transaction.save()
		}
	}

	return transaction
}

// Removes from accounts balance. Returns created transaction with details
func (tm *TransactionManager) MakeRemove(account *accountsbundle.Account, sum int64) *Transaction {
	accepted := account.Balance >= sum

	transaction := makeTransaction(TypeRemove, account, "", sum, accepted)

	if transaction != nil {
		if !account.DecreaseBalance(sum) {
			transaction.Accepted = false
			transaction.save()
		}
	}

	return transaction
}

// Returns all transactions from (including) to (excluding) dates
func (tm *TransactionManager) GetTransactions(fromYear, fromMonth, fromDay, toYear, toMonth, toDay int) *[]Transaction {
	timeStart := time.Date(fromYear, time.Month(fromMonth), fromDay, 0, 0, 0, 0, time.Local)
	start := timeStart.Format(time.RFC3339)

	timeEnd := time.Date(toYear, time.Month(toMonth), toDay, 0, 0, 0, 0, time.Local)
	end := timeEnd.Format(time.RFC3339)

	transactions, err := getTransactionsByTimeFromDB(start, end)
	if err != nil {
		log.Println(err)
		return nil
	}
	return transactions
}

// Returns all transactions for specific account, from (including) to (excluding) dates
func (tm *TransactionManager) GetTransactionsByAccount(account *accountsbundle.Account, fromYear, fromMonth, fromDay, toYear, toMonth, toDay int) *[]Transaction {
	timeStart := time.Date(fromYear, time.Month(fromMonth), fromDay, 0, 0, 0, 0, time.Local)
	start := timeStart.Format(time.RFC3339)

	timeEnd := time.Date(toYear, time.Month(toMonth), toDay, 0, 0, 0, 0, time.Local)
	end := timeEnd.Format(time.RFC3339)

	transactions, err := getTransactionsByAccountByTimeFromDB(*account, start, end)
	if err != nil {
		log.Println(err)
		return nil
	}
	return transactions
}

// END TransactionManager

// ************************************************************

// START Transaction

// Transaction struct, containing most information about a transaction
type Transaction struct {
	ID 			int				`json:"id" db:"id"`
	Type 		TransactionType	`json:"type" db:"type"`
	AccountID	int				`json:"account_id" db:"account_id"`
	CardID		string			`json:"card_id,omitempty" db:"card_id"`
	Sum			int64			`json:"sum" db:"sum"`
	Accepted 	bool			`json:"accepted" db:"accepted"`
	Time 		string			`json:"time" db:"time"`
}

// Updates transactions details in DB
func (trans *Transaction) save() bool {
	err := updateTransactionInDB(*trans)
	if err != nil {
		log.Println(err)
		return false
	}
	return false
}

// Makes a transaction and stores it in DB
func makeTransaction(t TransactionType, account *accountsbundle.Account, cardID string, sum int64, accepted bool) *Transaction {
	transaction := &Transaction{
		Type: t,
		AccountID: account.ID,
		CardID: cardID,
		Sum: sum,
		Accepted: accepted,
		Time: time.Now().Format(time.RFC3339),
	}

	err := addTransactionToDB(transaction)

	if err != nil {
		log.Println(err)
		return nil
	}
	return transaction
}

// END Transaction

// ************************************************************

// START TransactionType

// Defined int for defining transaction type
type TransactionType int

// Constant transaction type codes
const (
	TypePurchase 	TransactionType	= 1
	TypeInsert		TransactionType	= 2
	TypeRemove		TransactionType = 3
)

// END TransactionType