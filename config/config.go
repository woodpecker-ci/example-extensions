package config

import (
	"crypto/ed25519"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterConfigExtension(r *gin.RouterGroup, pubKey ed25519.PublicKey) {
	r.POST("/ciconfig", func(c *gin.Context) {
		var req incoming
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
