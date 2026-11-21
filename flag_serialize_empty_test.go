package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringSliceSetFromSerializedEmpty(t *testing.T) {
	var dest []string
	v := StringSlice{}.Create([]string{"x"}, &dest, StringConfig{})

	require.NoError(t, v.Set(NewStringSlice().Serialize()))
	assert.Nil(t, v.Get())
	assert.Nil(t, dest)

	require.NoError(t, v.Set("new"))
	assert.Equal(t, []string{"new"}, v.Get())
	assert.Equal(t, []string{"new"}, dest)

	// a non-empty serialized value still overwrites and stays connected
	require.NoError(t, v.Set(NewStringSlice("a", "b").Serialize()))
	assert.Equal(t, []string{"a", "b"}, v.Get())
	assert.Equal(t, []string{"a", "b"}, dest)
}

func TestStringMapSetFromSerializedEmpty(t *testing.T) {
	var dest map[string]string
	v := StringMap{}.Create(map[string]string{"x": "y"}, &dest, StringConfig{})

	require.NoError(t, v.Set(NewStringMap(nil).Serialize()))
	assert.Nil(t, v.Get())
	assert.Nil(t, dest)

	require.NoError(t, v.Set("new=val"))
	assert.Equal(t, map[string]string{"new": "val"}, v.Get())
	assert.Equal(t, map[string]string{"new": "val"}, dest)
}

func TestCollectionFlagsSerializedEmptyFromArgs(t *testing.T) {
	var sl []string
	var m map[string]string
	cmd := &Command{
		Name: "test",
		Flags: []Flag{
			&StringSliceFlag{Name: "s", Destination: &sl},
			&StringMapFlag{Name: "m", Destination: &m},
		},
		Action: func(context.Context, *Command) error { return nil },
	}
	err := cmd.Run(context.Background(), []string{
		"test",
		"--s", NewStringSlice().Serialize(), "--s", "new",
		"--m", NewStringMap(nil).Serialize(), "--m", "k=v",
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"new"}, sl)
	assert.Equal(t, []string{"new"}, cmd.StringSlice("s"))
	assert.Equal(t, map[string]string{"k": "v"}, m)
	assert.Equal(t, map[string]string{"k": "v"}, cmd.StringMap("m"))
}
