package accountsbundle

import (
	"time"
	"log"
)

// Card for account
// This is basically a payment method for Account.
type Card struct {
	ID			int 		`json:"id" db:"id"`
	CardID		string 		`json:"card_id" db:"card_id"`
	AccountID	int 		`json:"-" db:"account_id"`
	Account		*Account 	`json:"account" db:"account"`
	Disabled	bool		`json:"disabled" db:"disabled"`
	Added 		string		`json:"added" db:"added"`
}

// Returns all cards and the accounts associated with them.
func GetCards() *[]Card {
	cards, err := getCardsFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return cards
}

// Get a card and its parent account by the cards unique ID.
func GetCardByCardID(cardID string) *Card {
	card, err := getCardByCardIDFromDB(cardID)
	if err != nil {
		log.Println(err)
		return nil
	}
	return card
}

// Add card to account and save to DB.
func (a *Account) AddCard(cardID string) *Card {
	card := &Card{
		CardID: cardID,
		AccountID: a.ID,
		Account: a,
		Disabled: false,
		Added: time.Now().Format(time.RFC3339),
	}
	if err := addCardToDB(card); err != nil {
		log.Println(err)
		return nil
	}
	return card
}

// Get all cards associated with account.
func (a *Account) GetCards() *[]Card {
	cards, err := getCardsByAccountFromDB(a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return cards
}

// Delete a card from account. Removes it from the DB too.
func (a *Account) DeleteCard(card *Card) bool {
	if a.ID != card.AccountID {
		return false
	}

	if err := deleteCardFromDB(card.ID); err != nil {
		log.Println(err)
		return false
	}
	return true
}