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
	"github.com/osoderholm/eletab-lite/app/common"
)

type Authenticator struct {
	signKey		*rsa.PrivateKey
	VerifyKey 	*rsa.PublicKey
	common.Controller
}

const (
	keyPath 		= "keys"
	privateKeyFile 	= "eletab.rsa.pem"
	publicKeyFile 	= "eletab.rsa.pub.pem"
)

func Init() *Authenticator {
	createKeys()

	return &Authenticator{
		signKey: getPrivateKey(),
		VerifyKey: getPublicKey(),
	}
}

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

func getPrivateKey() *rsa.PrivateKey {
	privateBytes, err := ioutil.ReadFile(path.Join(keyPath, privateKeyFile))
	fatal(err)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	fatal(err)

	return privateKey
}

func getPublicKey() *rsa.PublicKey {
	publicBytes, err := ioutil.ReadFile(path.Join(keyPath, publicKeyFile))
	fatal(err)

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	fatal(err)

	return publicKey
}

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

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
