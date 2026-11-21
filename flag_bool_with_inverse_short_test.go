package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoolWithInverseFlag_ShortOptionHandling(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		cFlag    bool
	}{
		{name: "single alias", args: []string{"app", "-v"}, expected: true},
		{name: "combined last", args: []string{"app", "-cv"}, expected: true, cFlag: true},
		{name: "combined first", args: []string{"app", "-vc"}, expected: true, cFlag: true},
		{name: "not set", args: []string{"app", "-c"}, expected: false, cFlag: true},
		{name: "inverse long", args: []string{"app", "-c", "--no-verbose"}, expected: false, cFlag: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var verbose, c bool
			cmd := &Command{
				Name:                   "app",
				UseShortOptionHandling: true,
				Flags: []Flag{
					&BoolWithInverseFlag{Name: "verbose", Aliases: []string{"v"}, Destination: &verbose},
					&BoolFlag{Name: "c", Destination: &c},
				},
				Action: func(context.Context, *Command) error { return nil },
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.expected, verbose)
			assert.Equal(t, tc.expected, cmd.Bool("verbose"))
			assert.Equal(t, tc.cFlag, c)
		})
	}
}

func TestCommand_UseShortOptionHandling_combinedValueFlagName(t *testing.T) {
	var gotName string
	cmd := &Command{
		Name:                   "app",
		UseShortOptionHandling: true,
		Flags: []Flag{
			&BoolFlag{Name: "c"},
			&StringFlag{Name: "name", Aliases: []string{"n"}, Destination: &gotName},
		},
		Action: func(context.Context, *Command) error { return nil },
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-cn", "foo"}))
	assert.Equal(t, "foo", gotName)
	assert.True(t, cmd.Bool("c"))
	assert.Equal(t, "foo", cmd.String("n"))
}
