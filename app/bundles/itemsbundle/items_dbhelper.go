package itemsbundle

import (
	"github.com/osoderholm/eletab-lite/app/common"
	"log"
)

// Empty helper struct for DBHelper interface
type helper struct {}

// Mandatory table creation function
func (helper helper) CreateTable(database *common.Database) error {
	qryItem := `
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			price INT64 NOT NULL,
			category_id INTEGER NOT NULL,
			added DATETIME NOT NULL
			);`

	_, err := database.Exec(qryItem)

	if err != nil {
		log.Println(err)
		return err
	}

	qryCategory := `
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			created DATETIME NOT NULL
			);`

	_, err = database.Exec(qryCategory)

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

// START Item DB

func addItemToDB(item *Item) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `
		INSERT INTO items (
			name, price, category_id, added
		) VALUES (
			:name, :price, :category.id, :added
		);
		`
	res, err := db.NamedExec(query, item)

	if err != nil { return err }

	id, err := res.LastInsertId()

	if err != nil { return err }

	item.ID = int(id)

	return nil
}

func getItemsFromDB() (*[]Item, error) {
	var items []Item

	db, err := openDB()
	defer db.Close()
	if err != nil { return &items, err }

	query := `SELECT
				items.*,
				categories.id "category.id",
				categories.name "category.name",
				categories.created "category.created"
			FROM
				items JOIN categories ON items.category_id = categories.id;`

	return &items, db.Select(&items, query)
}

func getItemByIDFromDB(itemID int) (*Item, error) {
	var item Item

	db, err := openDB()
	defer db.Close()
	if err != nil { return &item, err }

	query := `SELECT
				items.*,
				categories.id "category.id",
				categories.name "category.name",
				categories.created "category.created"
			FROM
				items JOIN categories ON items.category_id = categories.id
			 WHERE items.id = ?;`

	return &item, db.Get(&item, query, itemID)
}

func getItemsByCategoryFromDB(category Category) (*[]Item, error) {
	var items []Item

	db, err := openDB()
	defer db.Close()
	if err != nil { return &items, err }

	query := `SELECT
				items.*,
				categories.id "category.id",
				categories.name "category.name",
				categories.created "category.created"
			FROM
				items JOIN categories ON items.category_id = categories.id
			WHERE items.category = ?;`

	return &items, db.Get(&items, query, category.ID)
}

func deleteItemFromDB(itemID int) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `DELETE FROM items WHERE items.id = ?;`

	_, err = db.Exec(query, itemID)
	return err
}

// END Item DB

// ********************************************************

// START Category DB

func addCategoryToDB(category *Category) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `
		INSERT INTO categories (
			name, created
		) VALUES (
			:name, :created
		);
		`
	res, err := db.NamedExec(query, category)

	if err != nil { return err }

	id, err := res.LastInsertId()

	if err != nil { return err }

	category.ID = int(id)

	return nil
}

func getCategoriesFromDB() (*[]Category, error) {
	var categories []Category

	db, err := openDB()
	defer db.Close()
	if err != nil { return &categories, err }

	query := `SELECT * FROM categories;`

	return &categories, db.Select(&categories, query)
}

func getCategoryFromDB(categoryID int) (*Category, error) {
	var category Category

	db, err := openDB()
	defer db.Close()
	if err != nil { return &category, err }

	query := `SELECT * FROM categories WHERE categories.id = ?;`

	return &category, db.Get(&category, query, categoryID)
}

func deleteCategoryFromDB(categoryID int) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `DELETE FROM categories WHERE categories.id = ?;`

	_, err = db.Exec(query, categoryID)
	return err
}

// END Category DB
