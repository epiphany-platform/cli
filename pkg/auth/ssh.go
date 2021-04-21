package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"path"

	"github.com/epiphany-platform/cli/internal/logger"

	"golang.org/x/crypto/ssh"
)

const rsaKeyName = "vms_rsa"

type RsaKeyPair struct {
	Name string
}

func GenerateRsaKeyPair(directory string) (RsaKeyPair, error) {
	err := generateRsaKeyPair(directory, rsaKeyName)
	if err != nil {
		logger.Error().Err(err).Msg("generation of RSA key pair failed")
		return RsaKeyPair{}, err
	}
	return RsaKeyPair{rsaKeyName}, nil
}

func generateRsaKeyPair(directory string, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.Error().Err(err).Msg("rsa.GenerateKey failed")
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		logger.Error().Err(err).Msg("ssh.NewPublicKey failed")
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		logger.Error().Err(err).Msgf("write to private key file %s failed", path.Join(directory, name))
		return err
	}
	err = ioutil.WriteFile(path.Join(directory, name+".pub"), publicKeyBytes, 0644)
	if err != nil {
		logger.Error().Err(err).Msgf("write to public key file %s failed", path.Join(directory, name+".pub"))
		return err
	}
	logger.Debug().Msgf("correctly saved private and public key files: %s and %s", path.Join(directory, name), path.Join(directory, name+".pub"))
	return nil
}
