package utils

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/yaronf/httpsign"
	"golang.org/x/oauth2"
)

func getPubKeyFromServer(url, token string) ([]byte, error) {
	ctx := context.Background()
	pubKeyUrl := fmt.Sprintf("%s/api/signature/public-key", url)

	config := new(oauth2.Config)
	client := config.Client(
		ctx,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	pubKeyResponse, err := client.Get(pubKeyUrl)
	if err != nil {
		return nil, errors.New("Failed to get public key file " + err.Error())
	}

	pubKeyRaw, err := io.ReadAll(pubKeyResponse.Body)
	if err != nil {
		return nil, errors.New("Failed to read public key file " + err.Error())
	}

	if len(pubKeyRaw) == 0 || string(pubKeyRaw) == "User not authorized" {
		return nil, errors.New("Failed to get public key file")
	}

	return pubKeyRaw, nil
}

func getPubKey() ([]byte, error) {
	woodpeckerServerURL := os.Getenv("EXTENSION_WOODPECKER_URL")
	woodpeckerToken := os.Getenv("EXTENSION_WOODPECKER_TOKEN")
	if woodpeckerServerURL != "" && woodpeckerToken != "" {
		return getPubKeyFromServer(woodpeckerServerURL, woodpeckerToken)
	}

	localFilePath := os.Getenv("EXTENSION_PUBLIC_KEY_FILE")
	if localFilePath != "" {
		pubKeyRaw, err := os.ReadFile(localFilePath)
		if err != nil {
			return nil, errors.New("Failed to read public key file " + err.Error())
		}

		return pubKeyRaw, nil
	}

	return nil, errors.New("EXTENSION_WOODPECKER_URL is not set")
}

func GetPubKey() (ed25519.PublicKey, error) {
	pubKeyRaw, err := getPubKey()
	if err != nil {
		return nil, err
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

// Verify check cryptographic signature
func Verify(pubKey ed25519.PublicKey, r *http.Request) error {
	pubKeyID := "woodpecker-ci-extensions"

	verifier, err := httpsign.NewEd25519Verifier(pubKey,
		httpsign.NewVerifyConfig(),
		httpsign.Headers("@request-target", "content-digest"))
	if err != nil {
		return err
	}

	err = httpsign.VerifyRequest(pubKeyID, *verifier, r)
	if err != nil {
		return err
	}

	return nil
}
