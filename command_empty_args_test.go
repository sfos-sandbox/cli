package cli

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand_RunWithEmptyArgs(t *testing.T) {
	tests := []struct {
		name string
		cmd  func(ran *bool) *Command
	}{
		{
			name: "root without name",
			cmd: func(ran *bool) *Command {
				return &Command{
					Writer: io.Discard,
					Action: func(context.Context, *Command) error { *ran = true; return nil },
				}
			},
		},
		{
			name: "shell completion enabled",
			cmd: func(ran *bool) *Command {
				return &Command{
					Name:                  "foo",
					EnableShellCompletion: true,
					Writer:                io.Discard,
					Action:                func(context.Context, *Command) error { *ran = true; return nil },
				}
			},
		},
		{
			name: "root without name and shell completion enabled",
			cmd: func(ran *bool) *Command {
				return &Command{
					EnableShellCompletion: true,
					Writer:                io.Discard,
					Action:                func(context.Context, *Command) error { *ran = true; return nil },
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, args := range [][]string{{}, nil} {
				ran := false
				cmd := tc.cmd(&ran)
				require.NotPanics(t, func() {
					assert.NoError(t, cmd.Run(context.Background(), args))
				})
				assert.True(t, ran, "action should have run")
			}
		})
	}
}

func TestCheckShellCompleteFlag_EmptyArgs(t *testing.T) {
	cmd := &Command{EnableShellCompletion: true}
	ok, args := checkShellCompleteFlag(cmd, []string{})
	assert.False(t, ok)
	assert.Empty(t, args)
}
