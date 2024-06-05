package config

import (
	"crypto/ed25519"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

func RegisterConfigExtension(r *gin.RouterGroup, pubKey ed25519.PublicKey) {
	filterRegex := os.Getenv("EXTENSION_CONFIG_REPO_FILTER")
	filter := regexp.MustCompile(filterRegex)

	r.POST("/ciconfig", func(c *gin.Context) {
		var req incoming
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !filter.MatchString(req.Repo.Name) {
			// use default config
			c.JSON(http.StatusNoContent, gin.H{"message": "No config override for this repo"})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{"configs": []config{
			{
				Name: "central pipe",
				Data: overrideConfiguration,
			},
		}})
	})
}
