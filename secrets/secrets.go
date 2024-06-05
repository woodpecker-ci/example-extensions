package secrets

import (
	"crypto/ed25519"
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.woodpecker-ci.org/woodpecker/v2/server/model"
)

func RegisterSecretsExtension(r *gin.RouterGroup, pubKey ed25519.PublicKey) {
	secrets := make([]*model.Secret, 0)

	r.GET("/secrets", func(c *gin.Context) {
		c.JSON(http.StatusOK, secrets)
	})

	r.GET("/secrets/:secretName", func(c *gin.Context) {
		for _, secret := range secrets {
			if secret.Name == c.Param("secretName") {
				c.JSON(http.StatusOK, secrets)
				break
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
	})

	r.POST("/secrets", func(c *gin.Context) {
		var secret model.Secret
		if err := c.BindJSON(&secret); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		secrets = append(secrets, &secret)
		c.JSON(http.StatusCreated, secret)
	})

	r.POST("/secrets/:secretName", func(c *gin.Context) {
		for _, secret := range secrets {
			if secret.Name == c.Param("secretName") {
				if err := c.BindJSON(&secret); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, secret)
				break
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
	})

	r.DELETE("/secrets/:secretName", func(c *gin.Context) {
		for i, secret := range secrets {
			if secret.Name == c.Param("secretName") {
				secrets = append(secrets[:i], secrets[i+1:]...)
				c.JSON(http.StatusNoContent, nil)
				break
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
	})
}
