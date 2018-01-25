package accountsbundle

import (
	"time"
	"log"
)

type Card struct {
	ID			int 		`json:"id" db:"id"`
	CardID		string 		`json:"card_id" db:"card_id"`
	AccountID	int 		`json:"-" db:"account_id"`
	Account		*Account 	`json:"account" db:"account"`
	Disabled	bool		`json:"disabled" db:"disabled"`
	Added 		string		`json:"added" db:"added"`
}

func GetCards() *[]Card {
	cards, err := getCardsFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return cards
}

func GetCardByCardID(cardID string) *Card {
	card, err := getCardByCardIDFromDB(cardID)
	if err != nil {
		log.Println(err)
		return nil
	}
	return card
}

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

func (a *Account) GetCards() *[]Card {
	cards, err := getCardsByAccountFromDB(a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return cards
}

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