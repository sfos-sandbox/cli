package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericFlag_SchemaItemsType_NoPanic(t *testing.T) {
	f := &GenericFlag{Name: "gen"}
	assert.NotPanics(t, func() {
		assert.Equal(t, "", f.SchemaItemsType())
	})

	f = &GenericFlag{Name: "gen", Value: &Parser{"a", "b"}}
	assert.NotPanics(t, func() {
		assert.Equal(t, "", f.SchemaItemsType())
	})
}

func TestGenericFlag_NilValueEmptyEnvVar_NoPanic(t *testing.T) {
	t.Setenv("APP_GEN_EMPTY", "")

	cmd := &Command{
		Name: "app",
		Flags: []Flag{
			&GenericFlag{Name: "gen", Sources: EnvVars("APP_GEN_EMPTY")},
		},
		Action: func(context.Context, *Command) error { return nil },
	}

	var err error
	assert.NotPanics(t, func() {
		err = cmd.Run(buildTestContext(t), []string{"app"})
	})
	require.NoError(t, err)
	assert.True(t, cmd.IsSet("gen"))
	assert.Nil(t, cmd.Generic("gen"))
}
