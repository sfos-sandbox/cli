package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand_PositionalArgsPreserveWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{name: "padded", args: []string{"prog", "  padded  "}, expected: []string{"  padded  "}},
		{name: "leading", args: []string{"prog", " lead"}, expected: []string{" lead"}},
		{name: "trailing", args: []string{"prog", "trail\t"}, expected: []string{"trail\t"}},
		{name: "mixed with flag", args: []string{"prog", "--foo", "bar", " a ", "b"}, expected: []string{" a ", "b"}},
		{name: "whitespace only", args: []string{"prog", "   "}, expected: []string{"   "}},
		{name: "after terminator", args: []string{"prog", "--", " x "}, expected: []string{" x "}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			cmd := &Command{
				Name:  "prog",
				Flags: []Flag{&StringFlag{Name: "foo"}},
				Action: func(_ context.Context, cmd *Command) error {
					got = cmd.Args().Slice()
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expected, cmd.Args().Slice())
		})
	}
}
