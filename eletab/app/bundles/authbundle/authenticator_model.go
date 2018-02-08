package authbundle

import (
	"os"
	"path"
	"log"
	"crypto/rsa"
	"crypto/rand"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"github.com/dgrijalva/jwt-go"
	"github.com/osoderholm/eletab-lite/eletab/app/common"
)

// Used for applying authentication for accessing APIs
type Authenticator struct {
	signKey		*rsa.PrivateKey
	VerifyKey 	*rsa.PublicKey
	common.Controller
}

// File names and path for generated RSA keys
const (
	keyPath 		= "keys"
	privateKeyFile 	= "eletab.rsa.pem"
	publicKeyFile 	= "eletab.rsa.pub.pem"
)

// Initializes an authenticator.
// Generates RSA public and private keys if needed.
func Init() *Authenticator {
	createKeys()

	return &Authenticator{
		signKey: getPrivateKey(),
		VerifyKey: getPublicKey(),
	}
}

// Generate RSA private and public keys if none exist
func createKeys() {
	if _, err := os.Stat(path.Join(keyPath)); os.IsNotExist(err) {
		fatal(os.MkdirAll(path.Join(keyPath), os.ModePerm))
	}

	if _, err := os.Stat(path.Join(keyPath, privateKeyFile)); os.IsNotExist(err) {
		createPrivateKey()
	}

	if _, err := os.Stat(path.Join(keyPath, publicKeyFile)); os.IsNotExist(err) {

	}

}

// Read private key from file
func getPrivateKey() *rsa.PrivateKey {
	privateBytes, err := ioutil.ReadFile(path.Join(keyPath, privateKeyFile))
	fatal(err)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	fatal(err)

	return privateKey
}

// Read public key from file
func getPublicKey() *rsa.PublicKey {
	publicBytes, err := ioutil.ReadFile(path.Join(keyPath, publicKeyFile))
	fatal(err)

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	fatal(err)

	return publicKey
}

// Generates and writes private key to file
func createPrivateKey() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	fatal(err)
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	err = ioutil.WriteFile(path.Join(keyPath, privateKeyFile), pemdata, os.ModePerm)
	createPublicKey(privateKey)
}

// Generates and writes public key to file
func createPublicKey(priv *rsa.PrivateKey) {
	PubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	fatal(err)
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PUBLIC KEY",
			Bytes: PubASN1,
		},
	)
	err = ioutil.WriteFile(path.Join(keyPath, publicKeyFile), pemdata, os.ModePerm)
}

// Used for error handling.
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
