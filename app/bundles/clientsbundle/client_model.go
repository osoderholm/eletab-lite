package clientsbundle

import (
	api "github.com/osoderholm/eletab-lite/app/bundles/apibundle"
	"crypto/sha1"
	"math/rand"
	"time"
	"fmt"
	"strconv"
	"log"
)

type Client struct {
	ID 			int 				`json:"id" db:"id"`
	Key 		string 				`json:"api_key" db:"key"`
	Secret		string				`json:"-" db:"secret"`
	Level		api.APIAccessLevel	`json:"api_level" db:"level"`
	Added 		string				`json:"added" db:"added"`
}

func AddClient(level api.APIAccessLevel) *Client {

	client := &Client{
		Key: generateHash(),
		Secret: generateHash(),
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

func GetClients() *[]Client {
	clients, err := getClientsFromDB()

	if err != nil {
		log.Println(err)
		return nil
	}

	return clients
}

func GetClientByKey(key string) *Client {
	client, err := getClientByKeyFromDB(key)

	if err != nil {
		log.Println(err)
		return nil
	}

	return client
}

func generateHash() string {
	rand.Seed(time.Now().UnixNano())
	h := sha1.New()

	h.Write([]byte(strconv.Itoa(rand.Int()) + time.Now().Format(time.RFC3339) + strconv.Itoa(rand.Int())))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)

}