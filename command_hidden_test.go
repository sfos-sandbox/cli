package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandIsHidden(t *testing.T) {
	var visited []string
	var hiddenSeen []string

	cmd := &Command{
		Name:                  "app",
		EnableShellCompletion: true,
		Commands: []*Command{
			{Name: "visible", Commands: []*Command{{Name: "leaf"}}},
			{
				Name:   "secret",
				Hidden: true,
				Commands: []*Command{
					{Name: "inner", Commands: []*Command{{Name: "deepest"}}},
				},
			},
			{
				Name: "tree",
				Action: func(_ context.Context, c *Command) error {
					return c.Root().Walk(func(sub *Command) error {
						if sub.Name == helpName {
							return nil
						}
						if sub.IsHidden() {
							hiddenSeen = append(hiddenSeen, sub.Name)
							return nil
						}
						visited = append(visited, sub.Name)
						return nil
					})
				},
			},
		},
	}

	require.NoError(t, cmd.Run(context.Background(), []string{"app", "tree"}))

	assert.ElementsMatch(t, []string{"secret", "inner", "deepest", "completion", "bash", "zsh", "fish", "pwsh"}, hiddenSeen)
	assert.Contains(t, visited, "app")
	assert.Contains(t, visited, "visible")
	assert.Contains(t, visited, "leaf")
	assert.Contains(t, visited, "tree")

	// The Hidden field keeps the value the developer set.
	secret := cmd.Command("secret")
	require.NotNil(t, secret)
	assert.True(t, secret.Hidden)
	assert.False(t, secret.Command("inner").Hidden)
	assert.False(t, secret.Command("inner").Command("deepest").Hidden)
}

func TestCommandIsHiddenStandalone(t *testing.T) {
	assert.False(t, (&Command{Name: "x"}).IsHidden())
	assert.True(t, (&Command{Name: "x", Hidden: true}).IsHidden())
}

func TestHiddenSubcommandStillRuns(t *testing.T) {
	ran := false
	cmd := &Command{
		Name: "app",
		Commands: []*Command{
			{
				Name:   "secret",
				Hidden: true,
				Commands: []*Command{
					{Name: "inner", Action: func(context.Context, *Command) error {
						ran = true
						return nil
					}},
				},
			},
		},
	}
	require.NoError(t, cmd.Run(context.Background(), []string{"app", "secret", "inner"}))
	assert.True(t, ran)
}

func TestHiddenCompletionHelpListsSubcommands(t *testing.T) {
	var out bytes.Buffer
	cmd := &Command{
		Name:                  "app",
		EnableShellCompletion: true,
		Writer:                &out,
	}
	require.NoError(t, cmd.Run(context.Background(), []string{"app", "completion", "--help"}))
	for _, name := range []string{"bash", "zsh", "fish", "pwsh"} {
		assert.Contains(t, out.String(), name)
	}
}

func TestIsHiddenFromWithinHiddenSubcommandAction(t *testing.T) {
	var innerHidden, visibleHidden bool
	cmd := &Command{
		Name: "app",
		Commands: []*Command{
			{
				Name:   "secret",
				Hidden: true,
				Commands: []*Command{
					{Name: "inner", Action: func(_ context.Context, c *Command) error {
						innerHidden = c.IsHidden()
						return nil
					}},
				},
			},
			{Name: "visible", Action: func(_ context.Context, c *Command) error {
				visibleHidden = c.IsHidden()
				return nil
			}},
		},
	}
	require.NoError(t, cmd.Run(context.Background(), []string{"app", "secret", "inner"}))
	assert.True(t, innerHidden)
	assert.False(t, cmd.Command("secret").Command("inner").Hidden)

	require.NoError(t, cmd.Run(context.Background(), []string{"app", "visible"}))
	assert.False(t, visibleHidden)
}
