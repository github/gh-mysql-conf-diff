package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserPasswordFromEnv(t *testing.T) {
	t.Setenv("MYSQL_USER", "username")
	t.Setenv("MYSQL_PASSWORD", "password")
	user, password, err := getMySQLUserInfo()
	require.NoError(t, err)
	require.Equal(t, "username", user)
	require.Equal(t, "password", password)
}

func TestBasicFlags(t *testing.T) {
	context, err := newInputContext().parseArgs(
		[]string{"my.cnf", "localhost:1000", "--watch-options=option1,option2"})
	require.NoError(t, err)
	require.Equal(t, "my.cnf", context.configPath)
	require.Equal(t, "localhost:1000", context.serverAndPort)
	expected := map[string]any{"option1": true, "option2": true}
	if !reflect.DeepEqual(context.optionKeysToWatch, expected) {
		t.Fatalf("expected configOptions to be %v, got %v", expected, context.optionKeysToWatch)
	}
}

func TestBasicFlagsWithoutExecute(t *testing.T) {
	// test that when --apply-changes is not used, executeFlag is set to be off
	context, err := newInputContext().parseArgs(
		[]string{"my.cnf", "localhost:1000"})
	require.NoError(t, err)
	require.Equal(t, false, context.applyTheChanges)
}

func TestHelp(t *testing.T) {
	_, err := newInputContext().parseArgs(
		[]string{"--help"})
	require.Error(t, err)
	require.ErrorIs(t, errHelpFlagIsSet, err)
}

func TestExecuteRequiresWatchedOptions(t *testing.T) {
	// Try running --apply-changes without --watch-options
	_, err := newInputContext().parseArgs(
		[]string{"my.cnf", "localhost:1000", "--apply-changes"})
	// assert user error: that --watch-options is required
	require.Error(t, err)
	require.Contains(t, err.Error(), "--watch-options required")
}
