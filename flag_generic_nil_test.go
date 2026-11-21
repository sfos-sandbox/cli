package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericFlag_SchemaItemsType_NilValue(t *testing.T) {
	f := &GenericFlag{Name: "gen"}
	assert.NotPanics(t, func() {
		assert.Equal(t, "", f.SchemaItemsType())
	})
	assert.NotPanics(t, func() {
		assert.Equal(t, "", f.SchemaType())
	})
}

func TestGenericFlag_EmptySourceValue_NilValue(t *testing.T) {
	t.Setenv("GEN_EMPTY_ENV", "")

	f := &GenericFlag{Name: "gen", Sources: EnvVars("GEN_EMPTY_ENV")}
	cmd := &Command{
		Name:  "app",
		Flags: []Flag{f},
		Action: func(context.Context, *Command) error {
			return nil
		},
	}

	assert.NotPanics(t, func() {
		require.NoError(t, cmd.Run(buildTestContext(t), []string{"app"}))
	})
	assert.True(t, f.IsSet())
}
