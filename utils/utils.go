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

	"github.com/go-ap/httpsig"
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

func Verify(pubKey ed25519.PublicKey, w http.ResponseWriter, r *http.Request) error {
	// check signature
	pubKeyID := "woodpecker-ci-plugins"

	keystore := httpsig.NewMemoryKeyStore()
	keystore.SetKey(pubKeyID, pubKey)

	verifier := httpsig.NewVerifier(keystore)
	verifier.SetRequiredHeaders([]string{"(request-target)", "date"})

	keyID, err := verifier.Verify(r)
	if err != nil {
		return err
	}

	if keyID != pubKeyID {
		return err
	}

	return nil
}
