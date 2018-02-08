package clientsbundle

import (
	"log"
	"github.com/osoderholm/eletab-lite/eletab/app/common"
)

// Empty helper struct for DBHelper interface
type helper struct {}

// Mandatory table creation function
func (helper helper) CreateTable(database *common.Database) error {
	qry := `
		CREATE TABLE IF NOT EXISTS clients (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			key TEXT UNIQUE NOT NULL,
			secret TEXT NOT NULL,
			description TEXT NOT NULL,
			level INTEGER NOT NULL,
			added DATETIME NOT NULL
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

func addClientToDB(client *Client) error {
	db, err := openDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	if err != nil { return err }

	query := `
		INSERT INTO clients (
			key, secret, description, level, added
		) VALUES (
			:key, :secret, :description, :level, :added
		);
		`
	res, err := db.NamedExec(query, client)

	if err != nil { return err }

	id, err := res.LastInsertId()

	if err != nil { return err }

	client.ID = int(id)

	return nil

}

func getClientsFromDB() (*[]Client, error) {
	var clients []Client

	db, err := openDB()
	defer db.Close()
	if err != nil { return &clients, err }

	query := `SELECT * FROM clients;`

	return &clients, db.Select(&clients, query)
}

func getClientByKeyFromDB(key string) (*Client, error) {
	var client Client

	db, err := openDB()
	defer db.Close()
	if err != nil { return &client, err }

	query := `SELECT * FROM clients WHERE clients.key = ?;`

	return &client, db.Get(&client, query, key)
}

func deleteClientFromDB(clientID int) error {
	db, err := openDB()
	defer db.Close()
	if err != nil { return err }

	query := `DELETE FROM clients WHERE clients.id = ?;`

	_, err = db.Exec(query, clientID)
	return err
}