package apibundle

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"golang.org/x/net/html"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/osoderholm/eletab-lite/eletab/app/common"
	"github.com/osoderholm/eletab-lite/eletab/app/crypt"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/clientsbundle"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/authbundle"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/transactionsbundle"
	"github.com/osoderholm/eletab-lite/eletab/app/bundles/accountsbundle"
)

// Controller for the API bundle
type APIController struct {
	common.Controller
	TM	*transactionsbundle.TransactionManager
}

// Creates a new API controller
func NewController() *APIController {
	return &APIController{TM: transactionsbundle.NewManager()}
}

// Handle API root
func (c *APIController) Handle(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	clientKey := r.Form.Get("client")
	if len(clientKey) == 0 {
		log.Println("No client key in request")
		return
	}
	client := clientsbundle.GetClientByKey(clientKey)
	if client == nil {
		log.Println("No such client found, unauthorized!")
		c.SendErrorsJSON(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}
	// TODO: Add API help or something

}

// START Card functions

// For parsing authentication credentials for clients.
type clientCredientials struct {
	Key 	string 		`json:"api_key"`
	Secret 	string 		`json:"secret"`
}

// Handle card functions.
// Action is defined within by 'gorilla' vars
func (c *APIController) HandleCard(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	switch vars["action"] {
	case "get":
		c.handleCardGet(w, r)
		return

	case "charge":
		c.handleCardCharge(w, r)
		return
	}
	c.sendError(w, http.StatusBadRequest)
}

func (c *APIController) handleCardGet(w http.ResponseWriter, r *http.Request) {
	if ok := authorizeClient(r, clientsbundle.LevelCheck); !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	cardID := r.FormValue("card")

	if cardID != "" {
		card := accountsbundle.GetCardByCardID(cardID)
		if card == nil {
			c.sendError(w, http.StatusNoContent)
			return
		}
		c.SendJSON(w, http.StatusOK, card)
	} else {
		cards := accountsbundle.GetCards()
		if cards == nil {
			c.sendError(w, http.StatusNoContent)
			return
		}
		c.SendJSON(w, http.StatusOK, cards)
	}

}

func (c *APIController) handleCardCharge(w http.ResponseWriter, r *http.Request) {
	if ok := authorizeClient(r, clientsbundle.LevelCharge); !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	cardID := r.FormValue("card")
	sumStr := r.FormValue("sum")

	if cardID != "" && sumStr != "" {
		sum, err := strconv.Atoi(sumStr)
		if err != nil {
			c.sendError(w, http.StatusBadGateway)
			return
		}
		card := accountsbundle.GetCardByCardID(cardID)
		if card == nil {
			c.sendError(w, http.StatusNoContent)
			return
		}
		transaction := c.TM.MakePurchase(card.Account, card.CardID, int64(sum))
		c.SendJSON(w, http.StatusOK, transaction)
	} else {
		c.sendError(w, http.StatusBadRequest)
	}

}

func (c *APIController) HandleClientLogin(w http.ResponseWriter, r *http.Request) {
	var creds clientCredientials

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

// Returns whether authorization is ok and the client is authorized or not.
func authorizeClient(r *http.Request, level clientsbundle.AccessLevel) bool {
	r.ParseForm()
	clientKey := r.Form.Get("client")
	if len(clientKey) == 0 {
		log.Println("No client key in request")
		return false
	}
	client := clientsbundle.GetClientByKey(clientKey)
	if client == nil {
		log.Println("No such client found, unauthorized!")
		return false
	}
	return client.Level >= clientsbundle.AccessLevel(level)
}

// END Card functions

// Start Account functions

// For parsing account authentication details.
type accountCredentials struct {
	Username	string 	`json:"username"`
	Password	string 	`json:"password"`
}

// Handle account actions. The actions are defined with 'gorilla' vars.
func (c *APIController) HandleAccount(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	switch vars["action"] {
	case "get":
		c.handleAccountGet(w, r)
		return

	case "get_all":
		c.handleAccountGetAll(w, r)
		return

	case "increase":
		c.handleAccountInsert(w, r)
		return

	case "decrease":
		c.handleAccountRemove(w, r)
		return

	case "delete":
		c.handleAccountDelete(w, r)
		return

	case "edit":
		c.handleAccountEdit(w, r)
		return

	case "new":
		c.handleAccountNew(w, r)
		return

	case "get_cards":
		c.handleAccountGetCards(w, r)
		return

	case "add_card":
		c.handleAccountAddCard(w, r)
		return

	case "delete_card":
		c.handleAccountDeleteCard(w, r)
		return
	}
	c.sendError(w, http.StatusBadRequest)
}

// Get account from DB and do necessary error handling.
// Returns possible account (nil if not found) and http status code.
func getAccount(r *http.Request) (*accountsbundle.Account, int) {
	if r.ParseForm() != nil {
		return nil, http.StatusInternalServerError
	}
	accountID := formValue(r, "id")
	accountUsername := formValue(r, "username")

	if accountID != "" || accountUsername != "" {
		var account *accountsbundle.Account
		if accountID != "" {
			id, err := strconv.Atoi(accountID)
			if err != nil {
				return nil, http.StatusBadRequest
			}
			account = accountsbundle.GetAccountByID(id)
		} else {
			account = accountsbundle.GetAccountByUsername(accountUsername)
		}
		if account == nil {
			return nil, http.StatusNoContent
		}
		return account, http.StatusOK
	}
	return nil, http.StatusBadRequest
}

func (c *APIController) handleAccountGet(w http.ResponseWriter, r *http.Request) {
	tAccount, ok := authorizeAccount(r, accountsbundle.LevelDefault)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	account, code := getAccount(r)
	if account == nil {
		account = tAccount
	}
	if account != nil {
		if tAccount.Level >= accountsbundle.LevelModerator || account.ID == tAccount.ID {
			c.SendJSON(w, http.StatusOK, account)
		} else {
			c.sendError(w, http.StatusUnauthorized)
		}
	} else {
		c.sendError(w, code)
	}

}

func (c *APIController) handleAccountGetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelModerator)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	accounts := accountsbundle.GetAccounts()
	if accounts != nil {
		c.SendJSON(w, http.StatusOK, accounts)
	} else {
		c.sendError(w, http.StatusNoContent)
	}
}

func (c *APIController) handleAccountInsert(w http.ResponseWriter, r *http.Request) {
	if _, ok := authorizeAccount(r, accountsbundle.LevelModerator); !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	account, _ := getAccount(r)
	sum := formValue(r, "sum")
	if account != nil && sum != "" {
		s, err := strconv.Atoi(sum)
		if err != nil {
			c.sendError(w, http.StatusBadRequest)
			return
		}
		tr := c.TM.MakeInsert(account, int64(s))
		c.SendJSON(w, http.StatusOK, tr)
	} else {
		c.sendError(w, http.StatusBadRequest)
	}

}

func (c *APIController) handleAccountRemove(w http.ResponseWriter, r *http.Request) {
	if _, ok := authorizeAccount(r, accountsbundle.LevelModerator); !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	account, _ := getAccount(r)
	sum := formValue(r, "sum")
	if account != nil && sum != "" {
		s, err := strconv.Atoi(sum)
		if err != nil {
			c.sendError(w, http.StatusBadRequest)
			return
		}
		tr := c.TM.MakeRemove(account, int64(s))
		c.SendJSON(w, http.StatusOK, tr)
	} else {
		c.sendError(w, http.StatusBadRequest)
	}

}

func (c *APIController) handleAccountDelete(w http.ResponseWriter, r *http.Request) {
	tAccount, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	account, _ := getAccount(r)
	if account != nil {
		if tAccount.ID != account.ID {
			res := account.Delete()
			c.SendJSON(w, http.StatusOK, res)
		} else {
			c.sendError(w, http.StatusUnauthorized)
		}
	} else {
		c.sendError(w, http.StatusBadRequest)
	}

}

func (c *APIController) handleAccountEdit(w http.ResponseWriter, r *http.Request) {
	tAccount, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}
	account, _ := getAccount(r)
	if account == nil {
		c.sendError(w, http.StatusBadRequest)
		return
	}
	if account.Level >= tAccount.Level {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	name := formValue(r, "name")
	newPass := formValue(r, "new_pass")

	disabled, err := strconv.ParseBool(formValue(r, "disabled"))
	if tAccount.Level >=accountsbundle. LevelAdmin && err != nil {
		log.Println(err)
		c.sendError(w, http.StatusBadRequest)
		return
	}

	level, err := strconv.Atoi(formValue(r, "level"))
	if err != nil {
		log.Println(err)
		c.sendError(w, http.StatusBadRequest)
		return
	}
	if accountsbundle.AccountLevel(level) >= tAccount.Level {
		c.sendError(w, http.StatusUnauthorized)
		return
	}

	if len(name) == 0 {
		name = account.Name
	}

	if len(newPass) > 0 {
		hash, err := crypt.Encrypt(newPass)
		if err != nil {
			c.sendError(w, http.StatusInternalServerError)
			return
		}
		newPass = hash
	} else {
		newPass = account.Password
	}

	account.Level = accountsbundle.AccountLevel(level)
	account.Disabled = disabled
	account.Name = name
	account.Password = newPass

	res := account.Save()

	c.SendJSON(w, http.StatusOK, res)

}

func (c *APIController) handleAccountNew(w http.ResponseWriter, r *http.Request) {
	tAccount, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}

	var errs []string

	name := formValue(r, "name")
	username := formValue(r, "username")
	pass := formValue(r, "pass")
	if len(name) == 0 {
		errs = append(errs, "Missing value for 'name'")
	}
	if len(username) == 0 {
		errs = append(errs, "Missing value for 'username'")
	}
	if len(pass) == 0 {
		errs = append(errs, "Missing value for 'pass'")
	}
	level, err := strconv.Atoi(formValue(r, "level"))
	if err != nil || (err == nil && (level > accountsbundle.LevelSuperAdmin || level < accountsbundle.LevelDefault)) {
		log.Println(err)
		errs = append(errs, "Invalid level")
	}
	balance, err := strconv.Atoi(formValue(r, "balance"))
	if err != nil || balance < 0 {
		log.Println(err)
		errs = append(errs, "Invalid balance")
	}
	if accountsbundle.AccountLevel(level) >= tAccount.Level {
		errs = append(errs, "Not authorized to set that account level")
	}

	if len(errs) != 0 {
		c.SendErrorsJSON(w, http.StatusBadRequest, errs...)
		return
	}

	account := accountsbundle.AddAccount(name, username, pass, int64(balance), accountsbundle.AccountLevel(level))
	if account == nil {
		c.SendErrorsJSON(w, http.StatusInternalServerError, "Could not add account")
	}

	c.SendJSON(w, http.StatusOK, account)

}

func (c *APIController) handleAccountGetCards(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}

	account, _ := getAccount(r)
	if account == nil {
		c.sendError(w, http.StatusBadRequest)
		return
	}
	cards := account.GetCards()
	if cards == nil {
		c.sendError(w, http.StatusNoContent)
		return
	}
	c.SendJSON(w, http.StatusOK, cards)

}

func (c *APIController) handleAccountAddCard(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}

	account, _ := getAccount(r)
	if account == nil {
		c.sendError(w, http.StatusBadRequest)
		return
	}

	var errs []string

	cardId := formValue(r, "card_id")
	if len(cardId) == 0 {
		errs = append(errs, "Missing value for 'card_id'")
	}

	if len(errs) != 0 {
		c.SendErrorsJSON(w, http.StatusBadRequest, errs...)
		return
	}

	card := account.AddCard(cardId)

	c.SendJSON(w, http.StatusOK, card)

}

func (c *APIController) handleAccountDeleteCard(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	if r.ParseForm() != nil {
		c.sendError(w, http.StatusInternalServerError)
		return
	}



	account, _ := getAccount(r)
	if account == nil {
		c.sendError(w, http.StatusBadRequest)
		return
	}

	cardId := formValue(r, "card_id")
	if len(cardId) == 0 {
		c.SendErrorsJSON(w, http.StatusBadRequest, "Missing value for 'card_id'")
		return
	}

	card := accountsbundle.GetCardByCardID(cardId)

	if card == nil {
		c.SendErrorsJSON(w, http.StatusBadRequest, "Card not found")
		return
	}

	res := account.DeleteCard(card)

	c.SendJSON(w, http.StatusOK, res)

}

func (c *APIController) HandleAccountLogin(w http.ResponseWriter, r *http.Request) {
	var creds accountCredentials

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		log.Println(err)
		c.SendErrorsJSON(w, http.StatusForbidden, "Error in request")
		return
	}

	account := accountsbundle.GetAccountByUsername(creds.Username)
	if account == nil {
		c.SendErrorsJSON(w, http.StatusForbidden, "Invalid credentials")
		return
	}

	if ok, err := crypt.Compare(account.Password, creds.Password); !ok {
		log.Println(err)
		c.SendErrorsJSON(w, http.StatusForbidden, "Invalid credentials")
		return
	}

	claims := make(jwt.MapClaims)
	claims["account"] = account.ID

	expiration := 24*100 	// 100 days

	token, err := authbundle.GenerateToken(expiration, claims)

	if err != nil {
		c.SendErrorsJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	c.SendJSON(w, http.StatusOK, token)
}


func (c *APIController) HandleClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["action"] {
	case "new":
		c.handleClientNew(w, r)
		return

	case "get":
		c.handleClientGet(w, r)
		return

	case "get_all":
		c.handleClientGetAll(w, r)
		return

	case "delete":
		c.handleClientDelete(w, r)
		return
	}
	c.sendError(w, http.StatusBadRequest)
}

func getClient(r *http.Request) (*clientsbundle.Client, int) {
	if r.ParseForm() != nil {
		return nil, http.StatusInternalServerError
	}
	clientKey := formValue(r, "key")

	if clientKey != "" {
		client := clientsbundle.GetClientByKey(clientKey)
		if client == nil {
			return nil, http.StatusNoContent
		}
		return client, http.StatusOK
	}
	return nil, http.StatusBadRequest
}

func (c *APIController) handleClientNew(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}

	var errs []string

	description := formValue(r, "description")
	levelStr := formValue(r, "level")
	if len(description) == 0 {
		errs = append(errs, "Missing value for 'description'")
	}
	level, err := strconv.Atoi(levelStr)
	if err != nil || ( err == nil && (level < int(clientsbundle.LevelCheck) || level > int(clientsbundle.LevelEdit))) {
		log.Println("wow", err)
		errs = append(errs, "Invalid level")
	}

	if len(errs) != 0 {
		c.SendErrorsJSON(w, http.StatusBadRequest, errs...)
		return
	}

	client := clientsbundle.AddClient(description, clientsbundle.AccessLevel(level))

	if client == nil {
		c.SendErrorsJSON(w, http.StatusInternalServerError, "Could not add client")
	}

	c.SendJSON(w, http.StatusOK, client)

}

func (c *APIController) handleClientGet(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	client, code := getClient(r)
	if client != nil {
		c.SendJSON(w, http.StatusOK, client)
	} else {
		c.sendError(w, code)
	}

}

func (c *APIController) handleClientGetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	clients := clientsbundle.GetClients()
	if clients != nil {
		c.SendJSON(w, http.StatusOK, clients)
	} else {
		c.sendError(w, http.StatusNoContent)
	}
}

func (c *APIController) handleClientDelete(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizeAccount(r, accountsbundle.LevelAdmin)
	if !ok {
		c.sendError(w, http.StatusUnauthorized)
		return
	}
	client, code := getClient(r)
	if client != nil {
		res := clientsbundle.DeleteClient(client)
		c.SendJSON(w, http.StatusOK, res)
	} else {
		c.sendError(w, code)
	}
}

// Return whether authentication is ok and the account is authorized to perform task.
func authorizeAccount(r *http.Request, level accountsbundle.AccountLevel) (*accountsbundle.Account, bool) {
	r.ParseForm()
	accountId := r.Form.Get("account")
	if len(accountId) == 0 {
		log.Println("No account id in request")
		return nil, false
	}
	id, err := strconv.Atoi(accountId)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	account := accountsbundle.GetAccountByID(id)
	if account == nil {
		log.Println("No such account found, unauthorized!")
		return nil, false
	}
	return account, account.Level >= level
}

// END Account functions

// Strips all kinds of nasty injections from request form value and returns it.
func formValue(r *http.Request, key string) string {
	return html.EscapeString(strings.TrimSpace(r.FormValue(key)))
}

// Send a general error to client. If you need to specify a n error message, don't use this.
func (c *APIController) sendError(w http.ResponseWriter, code int) {
	errorMsg := "General error"
	switch code {
	case http.StatusInternalServerError:
		errorMsg = "Internal server error"
		break
	case http.StatusUnauthorized:
		errorMsg = "Unauthorized!"
		break
	case http.StatusBadRequest:
		errorMsg = "Bad request"
		break
	case http.StatusNoContent:
		errorMsg = "Not found"
		break
	default:
		code = http.StatusNoContent
		break
	}
	c.SendErrorsJSON(w, code, errorMsg)
}