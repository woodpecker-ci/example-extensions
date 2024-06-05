package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/woodpecker-ci/example-extensions/utils"
)

func main() {
	log.Println("Woodpecker sample secret extension")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	pubKey, err := utils.GetPubKey()
	if err != nil {
		log.Fatalf("Error getting public key: %v", err)
	}

	http.HandleFunc("/ciconfig", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := utils.Verify(pubKey, w, r)
		if err != nil {
			return
		}

		// var req incoming
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Failed to parse JSON"+err.Error(), http.StatusBadRequest)
			return
		}

		// TODO
	})

	err = http.ListenAndServe(os.Getenv("CONFIG_SERVICE_HOST"), nil)
	if err != nil {
		log.Fatalf("Error on listen: %v", err)
	}
}
