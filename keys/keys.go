package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/pmoieni/auth/log"
)

func GenerateKeys() {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Logger.Fatal("Cannot generate RSA key\n")
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("auth_pvt.pem")
	if err != nil {
		log.Logger.Fatal("error when create auth_pvt.pem: " + err.Error() + "\n")
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		log.Logger.Fatal("error when encode auth_pvt.pem: " + err.Error() + "\n")
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		log.Logger.Fatal("error when dumping publickey: " + err.Error() + "\n")
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("auth_pub.pem")
	if err != nil {
		log.Logger.Fatal("error when create auth_pub.pem: " + err.Error() + "\n")
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		log.Logger.Fatal("error when encode auth_pub pem: " + err.Error() + "\n")
	}
}
