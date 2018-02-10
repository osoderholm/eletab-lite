package authbundle

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
	"fmt"
)

// Token for authentication
type Token struct {
	Token 	string	`json:"token"`
}

// Function for applying client authentication middleware to request.
// The HandlerFunk is only called if authentication is successful
func (a *Authenticator) Handle(handler http.HandlerFunc) http.Handler {
	return negroni.New(
		negroni.HandlerFunc(a.ClientTokenValidatorMiddleware),
		negroni.Wrap(handler))
}

// Generates an authentication token. Takes number of hours until expiry
// and claims that are to be added to the JWT token.
func GenerateToken(expirationHours int, claims jwt.MapClaims) (Token, error) {
	a := Init()
	defer func() {a = nil}()

	token := jwt.New(jwt.SigningMethodRS256)

	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix()
	claims["iat"] = time.Now().Unix()

	token.Claims = claims

	var rToken Token

	tokenString, err := token.SignedString(a.signKey)

	if err != nil {
		return rToken, err
	}

	return Token{Token:tokenString}, err

}


// Used for refreshing the token, for basically infinite expiration date.
// Finds token in request and generates a new one with its claims.
func RefreshToken(r *http.Request, expirationHours int) (Token, error) {
	a := Init()
	defer func() {a = nil}()

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return a.VerifyKey, nil
		})
	if err != nil {
		log.Println("Unauthorized with error", err)
		return Token{}, err
	}
	if !token.Valid {
		log.Println("Unauthorized with invalid token")
		return Token{}, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return GenerateToken(expirationHours, claims)
}

// Middleware for validating JWT authentication tokens... duh.
// If validated, passes relevant token claim key-value pairs to next
// handler via http request.
func (a *Authenticator) ClientTokenValidatorMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)  {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return a.VerifyKey, nil
	})

	if err != nil {
		a.SendErrorsJSON(w, http.StatusUnauthorized, "Unauthorized!")
		log.Println("Unauthorized with error", err)
		return
	}
	if !token.Valid {
		a.SendErrorsJSON(w, http.StatusUnauthorized, "Invalid token!")
		log.Println("Unauthorized with invalid token")
		return
	}
	if r.ParseForm() != nil {
		a.SendErrorsJSON(w, http.StatusInternalServerError, "Internal server error")
		log.Println("Error parsing form")
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	for k, v := range claims {
		r.Form.Add(k, fmt.Sprintf("%v", v))
	}

	next(w, r)

}
