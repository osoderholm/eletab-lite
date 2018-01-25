package itemsbundle

import (
	"log"
	"time"
)

type Item struct {
	ID 			int 		`json:"id" db:"id"`
	Name 		string 		`json:"name" db:"name"`
	Price 		int64 		`json:"price" db:"price"`
	CategoryID 	int 		`json:"-" db:"category_id"`
	Category				`json:"category" db:"category"`
	Added 		string		`json:"-" db:"added"`
}

func AddItem(name string, price int64, category *Category) *Item {
	item := &Item{
		Name: name,
		Price: price,
		CategoryID: category.ID,
		Category: *category,
		Added: time.Now().Format(time.RFC3339),
	}
	if err := addItemToDB(item); err != nil {
		log.Println(err)
		return nil
	}
	return item
}

func GetItems() *[]Item {
	items, err := getItemsFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return items
}

func GetItemsByCategory(category *Category) *[]Item {
	items, err := getItemsByCategoryFromDB(*category)
	if err != nil {
		log.Println(err)
		return nil
	}
	return items
}

func GetItemByID(id int) *Item {
	item, err := getItemByIDFromDB(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return item
}
