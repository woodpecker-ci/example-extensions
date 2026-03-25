package config

import (
	_ "embed"

	"github.com/woodpecker-ci/example-extensions/types"
)

type config struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type incoming struct {
	*types.IncomingRequest
	
	Configuration []*config       `json:"configs"`
}

//go:embed central-pipeline-config.yaml
var overrideConfiguration string
