package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand_Run_NoArgs(t *testing.T) {
	tests := []struct {
		name string
		cmd  func(called *bool) *Command
	}{
		{
			name: "root command without name",
			cmd: func(called *bool) *Command {
				return &Command{
					Writer: &bytes.Buffer{},
					Action: func(context.Context, *Command) error {
						*called = true
						return nil
					},
				}
			},
		},
		{
			name: "shell completion enabled",
			cmd: func(called *bool) *Command {
				return &Command{
					Name:                  "foo",
					EnableShellCompletion: true,
					Writer:                &bytes.Buffer{},
					Action: func(context.Context, *Command) error {
						*called = true
						return nil
					},
				}
			},
		},
		{
			name: "root command without name and shell completion enabled",
			cmd: func(called *bool) *Command {
				return &Command{
					EnableShellCompletion: true,
					Writer:                &bytes.Buffer{},
					Action: func(context.Context, *Command) error {
						*called = true
						return nil
					},
				}
			},
		},
	}

	for _, tt := range tests {
		for _, args := range [][]string{nil, {}} {
			t.Run(tt.name, func(t *testing.T) {
				called := false
				cmd := tt.cmd(&called)
				require.NotPanics(t, func() {
					require.NoError(t, cmd.Run(buildTestContext(t), args))
				})
				require.True(t, called, "action should have been run")
			})
		}
	}
}

func TestCheckShellCompleteFlag_EmptyArgs(t *testing.T) {
	cmd := &Command{EnableShellCompletion: true}
	ok, args := checkShellCompleteFlag(cmd, []string{})
	require.False(t, ok)
	require.Empty(t, args)
}
