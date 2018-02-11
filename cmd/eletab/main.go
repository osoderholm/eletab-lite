package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/jamiealquiza/envy"

	"github.com/gorilla/mux"
	"github.com/osoderholm/eletab-lite/bundles/apibundle"
	"github.com/osoderholm/eletab-lite/bundles/authbundle"
)

func main() {
	var port = flag.Int("port", 8080, "port ")
	var address = flag.String("address", "127.0.0.1", "address to listen on")
	var staticfiles = flag.String("staticfiles", "../../", "path to all static files") // TODO: Include these in the binary
	envy.Parse("ELETAB")                                                               // looks for ELETAB_PORT and ELETAB_STATICFILES
	flag.Parse()
	if *port <= 0 {
		panic("port is fuked")
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

	sp := path.Join(*staticfiles, "/static/")
	log.Printf("Static path: %s\n", sp)
	staticFileDirectory := http.Dir(sp)
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler)

	http.Handle("/", r)
	sck := fmt.Sprintf("%s:%d", *address, *port)
	fmt.Printf("We are up and running at %s!\n", sck)
	log.Fatal(http.ListenAndServe(sck, nil))

}
