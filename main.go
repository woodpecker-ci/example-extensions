package main

import (
	"log"
	"os"

	"github.com/woodpecker-ci/example-extensions/config"
	"github.com/woodpecker-ci/example-extensions/secrets"
	"github.com/woodpecker-ci/example-extensions/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Woodpecker sample extensions server")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. Copy '.env.example' to '.env': %v", err)
	}

	pubKey, err := utils.GetPubKey()
	if err != nil {
		log.Fatalf("Error getting public key: %v", err)
	}

	r := gin.Default()

	// IMPORTANT: We verify all incoming requests to our api to ensure they are signed and coming from Woodpecker
	r.Use(func(c *gin.Context) {
		err := utils.Verify(pubKey, c.Writer, c.Request)
		if err != nil {
			log.Printf("Failed to verify request: %v", err)
			c.JSON(401, gin.H{"error": "Failed to verify request"})
			c.Abort()
			return
		}
	})

	g := r.Group("/repo/:repoId")
	config.RegisterConfigExtension(g, pubKey)
	secrets.RegisterSecretsExtension(g, pubKey)

	host := os.Getenv("EXTENSION_HOST")
	if host != "" {
		log.Printf("Starting server on %s\n", host)
	}

	r.Run(host)
}
