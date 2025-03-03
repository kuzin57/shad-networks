package main

import (
	"testing"

	"github.com/kuzin57/shad-networks/cmd/app"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestValidateApp(t *testing.T) {
	conf := "../../config/config.yaml"

	err := fx.ValidateApp(app.Create(conf))
	require.NoError(t, err)
}
