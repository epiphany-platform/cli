package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"

	"golang.org/x/crypto/ssh"
)

const rsaKeyName = "vms_rsa"

type RsaKeyPair struct {
	Name string
}

func GenerateRsaKeyPair(directory string) (kp RsaKeyPair, err error) {
	err = generateRsaKeyPair(directory, rsaKeyName)
	if err != nil {
		return
	}
	return RsaKeyPair{rsaKeyName}, nil
}

func generateRsaKeyPair(directory string, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
