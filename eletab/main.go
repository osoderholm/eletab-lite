package main

import (
	"github.com/gorilla/mux"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/apibundle"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/authbundle"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/clientsbundle"
	"fmt"
	"net/http"
	"log"
)

func main() {
	/*category := itemsbundle.AddCategory("test category")
	fmt.Println(*category)

	item := itemsbundle.AddItem("test", 100, category)
	fmt.Println(*item)

	items := itemsbundle.GetItems()
	for _, it := range *items {
		fmt.Println(it)
	}

	card := account.AddCard("1234567890")
	fmt.Println(*card)

	j, _ := json.Marshal(accountsbundle.GetCardByCardID("1234567890"))

	fmt.Printf("%s", j)*/

	/*account := accountsbundle.GetAccountByUsername("odoo")
	fmt.Println(*account)

	tm := transactionsbundle.NewManager()
	//transaction := tm.MakeInsert(account, 600)
	//fmt.Println(*transaction)

	transactions := tm.GetTransactions(2018, 1, 20, 2018, 1, 21)
	for _, t := range *transactions {
		fmt.Println(t)
	}
	transactions = tm.GetTransactionsByAccount(account, 2018, 1, 20, 2018, 1, 21)
	for _, t := range *transactions {
		fmt.Println(t)
	}*/

	clients := clientsbundle.GetClients()
	for _, c := range *clients {
		fmt.Println(c)
	}

	a := authbundle.Init()

	r := mux.NewRouter()

	apiSR := r.PathPrefix("/api/v1/").Subrouter()

	apiCtrl := apibundle.NewController()

	apiSR.Handle("/", a.Handle(apiCtrl.Handle))
	apiSR.Handle("/cards/{action}", a.Handle(apiCtrl.HandleCard))
	apiSR.Handle("/account/{action}", a.Handle(apiCtrl.HandleAccount))
	apiSR.Handle("/clients/{action}", a.Handle(apiCtrl.HandleClient))
	apiSR.HandleFunc("/client_login", apiCtrl.HandleClientLogin).Methods(http.MethodPost)
	apiSR.HandleFunc("/account_login", apiCtrl.HandleAccountLogin).Methods(http.MethodPost)

	staticFileDirectory := http.Dir("./app/static/")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
