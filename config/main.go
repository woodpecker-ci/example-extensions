package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/woodpecker-ci/example-extensions/utils"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
)

type config struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type incoming struct {
	Repo          *model.Repo     `json:"repo"`
	Build         *model.Pipeline `json:"pipeline"`
	Configuration []*config       `json:"configs"`
}

//go:embed central-pipeline-config.yaml
var overrideConfiguration string

func main() {
	log.Println("Woodpecker central config server")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("CONFIG_SERVICE_HOST")
	filterRegex := os.Getenv("CONFIG_SERVICE_OVERRIDE_FILTER")

	if host == "" {
		log.Fatal("Please make sure CONFIG_SERVICE_HOST and CONFIG_SERVICE_PUBLIC_KEY_FILE are set properly")
	}

	pubKey, err := utils.GetPubKey()
	if err != nil {
		log.Fatalf("Error getting public key: %v", err)
	}

	filter := regexp.MustCompile(filterRegex)

	http.HandleFunc("/ciconfig", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := utils.Verify(pubKey, w, r)
		if err != nil {
			return
		}

		var req incoming
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

		if filter.MatchString(req.Repo.Name) {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(map[string]interface{}{"configs": []config{
				{
					Name: "central pipe",
					Data: overrideConfiguration,
				},
			}})
			if err != nil {
				log.Printf("Error on encoding json %v\n", err)
			}
		} else {
			w.WriteHeader(http.StatusNoContent) // use default config
			// No need to write a response body
		}
	})

	err = http.ListenAndServe(os.Getenv("CONFIG_SERVICE_HOST"), nil)
	if err != nil {
		log.Fatalf("Error on listen: %v", err)
	}
}
