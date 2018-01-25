package authbundle

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Token struct {
	Token 	string	`json:"token"`
}

func (a *Authenticator) Handle(handler http.HandlerFunc) http.Handler {
	return negroni.New(
		negroni.HandlerFunc(a.TokenValidatorMiddleware),
		negroni.Wrap(handler))
}

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

func (a *Authenticator) TokenValidatorMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)  {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return a.VerifyKey, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	log.Printf("Token claims: %v\n\r", claims)

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

	next(w, r)

}
