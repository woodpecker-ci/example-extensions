package config

import (
	_ "embed"

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
