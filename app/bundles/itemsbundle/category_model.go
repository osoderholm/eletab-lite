package itemsbundle

import (
	"time"
	"log"
)

type Category struct {
	ID			int			`json:"id" db:"id"`
	Name 		string		`json:"name" db:"name"`
	Created 	string		`json:"-" db:"created"`
}

// Adds category to database and returns created category
func AddCategory(name string) *Category {
	category := &Category{
		Name: name,
		Created: time.Now().Format(time.RFC3339),
	}

	if err := addCategoryToDB(category); err != nil {
		log.Println(err)
		return nil
	}

	return category
}

func GetCategories() *[]Category {
	categories, err := getCategoriesFromDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	return categories
}


