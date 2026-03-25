package types

import (
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

type IncomingRequest struct {
	Repo          *model.Repo     `json:"repo"`
	Build         *model.Pipeline `json:"pipeline"`
	Netrc         *model.Netrc    `json:"netrc"`
}
