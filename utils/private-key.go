package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	privateKey *rsa.PrivateKey
)

// GetPrivateKey : Lấy giá trị private key đọc từ config file
func GetPrivateKey() *rsa.PrivateKey {
	if privateKey == nil {
		log.Println("Init private key ...")

		var folder string
		env := os.Getenv("APPLICATION_ENV")

		switch env {
		case "master", "dev", "uat", "sandbox":
			folder = env
		default:
			folder = "dev"
		}

		fileName := "config/" + folder + "/private.key"
		privateKeyBytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(fmt.Sprintf("Not found private key at %s", fileName))
		}
		privateKey, err = parseRsaPrivateKeyFromPemStr(string(privateKeyBytes))
		if err != nil {
			panic("Private key invalid format")
		}
	}

	return privateKey
}

func parseRsaPrivateKeyFromPemStr(privateKeyPEM string) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	return
}
