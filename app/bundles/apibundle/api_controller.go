package apibundle

import (
	"github.com/osoderholm/eletab-lite/app/common"
	"net/http"
	"encoding/json"
	"log"
	"github.com/osoderholm/eletab-lite/app/bundles/clientsbundle"
	"github.com/osoderholm/eletab-lite/app/bundles/authbundle"
	"github.com/dgrijalva/jwt-go"
)

// Controller for the API bundle
type APIController struct {
	common.Controller
}

type apiCredientials struct {
	Key 	string 		`json:"api_key"`
	Secret 	string 		`json:"secret"`
}

// Creates a new API controller
func NewController() *APIController {
	return &APIController{}
}

func (c *APIController) Handle(w http.ResponseWriter, r *http.Request)  {
	
}

func (c *APIController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds apiCredientials

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		log.Println(err)
		c.SendErrorsJSON(w, http.StatusForbidden, "Error in request")
		return
	}

	client := clientsbundle.GetClientByKey(creds.Key)
	if client == nil || client.Secret != creds.Secret {
		c.SendErrorsJSON(w, http.StatusForbidden, "Invalid credentials")
		return
	}

	claims := make(jwt.MapClaims)
	claims["client"] = client.Key

	token, err := authbundle.GenerateToken(24, claims)

	if err != nil {
		c.SendErrorsJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	c.SendJSON(w, http.StatusOK, token)
}


