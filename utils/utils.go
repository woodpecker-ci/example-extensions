package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-ap/httpsig"
)

func GetPubKey() (ed25519.PublicKey, error) {
	pubKeyPath := os.Getenv("CONFIG_SERVICE_PUBLIC_KEY_FILE") // Key in format of the one fetched from http(s)://your-woodpecker-server/api/signature/public-key
	if pubKeyPath == "" {
		return nil, errors.New("Please make sure CONFIG_SERVICE_HOST and CONFIG_SERVICE_PUBLIC_KEY_FILE are set properly")
	}

	pubKeyRaw, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, errors.New("Failed to read public key file")
	}

	pemblock, _ := pem.Decode(pubKeyRaw)

	b, err := x509.ParsePKIXPublicKey(pemblock.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse public key file " + err.Error())
	}
	pubKey, ok := b.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("Failed to parse public key file")
	}

	return pubKey, nil
}

func Verify(pubKey ed25519.PublicKey, w http.ResponseWriter, r *http.Request) error {
	// check signature
	pubKeyID := "woodpecker-ci-plugins"

	keystore := httpsig.NewMemoryKeyStore()
	keystore.SetKey(pubKeyID, pubKey)

	verifier := httpsig.NewVerifier(keystore)
	verifier.SetRequiredHeaders([]string{"(request-target)", "date"})

	keyID, err := verifier.Verify(r)
	if err != nil {
		log.Printf("config: invalid or missing signature in http.Request")
		http.Error(w, "Invalid or Missing Signature", http.StatusBadRequest)
		return err
	}

	if keyID != pubKeyID {
		log.Printf("config: invalid signature in http.Request")
		http.Error(w, "Invalid Signature", http.StatusBadRequest)
		return err
	}

	return nil
}
