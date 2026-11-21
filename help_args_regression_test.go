package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpFlagWithPositionalArgs(t *testing.T) {
	buildCmd := func(out *bytes.Buffer) *Command {
		return &Command{
			Name:   "app",
			Usage:  "root usage text",
			Writer: out,
			Commands: []*Command{
				{
					Name:      "command",
					Usage:     "command usage text",
					Arguments: []Argument{&StringArg{Name: "arg1"}},
					Action: func(context.Context, *Command) error {
						t.Fatal("action should not run when help is requested")
						return nil
					},
				},
				{
					Name:  "other",
					Usage: "other usage text",
				},
			},
		}
	}

	tests := []struct {
		name        string
		args        []string
		contains    string
		notContains string
	}{
		{"subcommand help flag", []string{"app", "command", "--help"}, "command usage text", "root usage text"},
		{"subcommand arg before help flag", []string{"app", "command", "ciao", "--help"}, "command usage text", "root usage text"},
		{"subcommand arg after help flag", []string{"app", "command", "--help", "ciao"}, "command usage text", "root usage text"},
		{"subcommand arg before short help flag", []string{"app", "command", "ciao", "-h"}, "command usage text", "root usage text"},
		{"root help flag with unknown arg", []string{"app", "--help", "ciao"}, "root usage text", ""},
		{"root unknown arg with short help flag", []string{"app", "ciao", "-h"}, "root usage text", ""},
		{"root help flag with subcommand name", []string{"app", "--help", "other"}, "other usage text", "root usage text"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := buildCmd(out)
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			require.Contains(t, out.String(), tc.contains)
			if tc.notContains != "" {
				require.NotContains(t, out.String(), tc.notContains)
			}
		})
	}
}
