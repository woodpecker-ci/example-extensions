package registries

import (
	"crypto/ed25519"
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"

	"github.com/woodpecker-ci/example-extensions/types"
)

//go:embed external-registries.json
var externalRegistries []byte

func RegisterRegistriesExtension(r *gin.RouterGroup, pubKey ed25519.PublicKey) {
	var registries []*model.Registry
	_ = json.Unmarshal(externalRegistries, &registries)

	r.POST("/registries", func(c *gin.Context) {
		var req types.IncomingRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, map[string]any{"registries": registries})
	})
}
