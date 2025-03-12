package main

import (
	"testing"

	"github.com/kuzin57/shad-networks/services/graph/cmd/app"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestValidateApp(t *testing.T) {
	var (
		conf    = "../../config/config.yaml"
		secrets = "../../config/secrets.yaml"
	)

	err := fx.ValidateApp(app.Create(conf, secrets))
	require.NoError(t, err)
}
