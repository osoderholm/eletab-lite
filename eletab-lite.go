package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/osoderholm/eletab-lite/bundles/apibundle"
	"github.com/osoderholm/eletab-lite/bundles/authbundle"
)

func main() {
	port := string(os.Getenv("ELETAB_PORT"))
	appPath := string(os.Getenv("ELETAB_PATH"))

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

	log.Println(path.Join(appPath, "/static/"))
	staticFileDirectory := http.Dir(path.Join(appPath, "/static/"))
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler)

	http.Handle("/", r)
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
