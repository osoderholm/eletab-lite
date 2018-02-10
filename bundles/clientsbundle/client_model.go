package clientsbundle

import (
	"crypto/sha1"
	"math/rand"
	"time"
	"fmt"
	"strconv"
	"log"
)

// Structure for a client that is able to authenticate
// with the server and interface with the API.
// It is worth noting that the secret is in plaintext, use with care.
type Client struct {
	ID 			int 				`json:"id" db:"id"`
	Key 		string 				`json:"api_key" db:"key"`
	Secret		string				`json:"secret" db:"secret"`
	Description	string				`json:"description" db:"description"`
	Level		AccessLevel			`json:"api_level" db:"level"`
	Added 		string				`json:"added" db:"added"`
}

// API Access levels, please don't mess with these
// unless you have a great transition plan for the upgrade
type AccessLevel int
const (
	LevelCheck		AccessLevel	= 1	// Can check account information
	LevelCharge		AccessLevel	= 2	// Can check account information and make purchases
	LevelEdit	 	AccessLevel	= 3	// Can edit, add and remove all account related information
)

// Creates a new client, adds it to the DB and returns it.
// APIAccessLevel documentation found in the apibundle package.
func AddClient(description string, level AccessLevel) *Client {

	client := &Client{
		Key: generateHash(),
		Secret: generateHash(),
		Description: description,
		Level: level,
		Added: time.Now().Format(time.RFC3339),
	}

	err := addClientToDB(client)

	if err != nil {
		log.Println(err)
		return nil
	}

	return client
}

// Returns all clients from DB
func GetClients() *[]Client {
	clients, err := getClientsFromDB()

	if err != nil {
		log.Println(err)
		return nil
	}

	return clients
}

// Get a specific client from DB by its key.
// Return nil if not found.
func GetClientByKey(key string) *Client {
	client, err := getClientByKeyFromDB(key)

	if err != nil {
		log.Println(err)
		return nil
	}

	return client
}

// Deletes client from DB, true if ok
func DeleteClient(client *Client) bool {
	return deleteClientFromDB(client.ID) == nil
}

// Generates a SHA1 hash and returns it as a string
// in hexadecimal form
func generateHash() string {
	rand.Seed(time.Now().UnixNano())
	h := sha1.New()

	h.Write([]byte(strconv.Itoa(rand.Int()) + time.Now().Format(time.RFC3339) + strconv.Itoa(rand.Int())))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)

}