package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockdbConn struct {
	mock.Mock
}

func TestMysqlConfDiff_DiffNoApply(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()
	confOptions := map[string]any{"key1": "value1"}
	serverVariables := map[string]any{"key1": "value2"}
	applyTheChanges := false

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges,
		&stdout, &stderr)

	// Check results
	expectedInStdOut := []string{"Difference found"}
	for _, expectedStr := range expectedInStdOut {
		assert.True(t, strings.Contains(stdout.String(), expectedStr))
	}
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestMysqlConfDiff_DiffApply(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()

	confOptions := map[string]any{"key1": "value1"}
	serverVariables := map[string]any{"key1": "value2"}
	applyTheChanges := true

	// Set SQL expectation (assuming applySetting executes 'SET key = value' query)
	m.ExpectExec("SET GLOBAL key1 = ?").WithArgs("value1").WillReturnResult(sqlmock.NewResult(1, 1))

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges,
		&stdout, &stderr)

	// Check results
	expectedInStdOut := []string{"Difference found", "Set variable", "key1 = value1"}
	for _, expectedStr := range expectedInStdOut {
		assert.True(t, strings.Contains(stdout.String(), expectedStr))
	}
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestMysqlConfDiff_NoDiff(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()

	confOptions := map[string]any{"key1": "value1"}
	serverVariables := map[string]any{"key1": "value1"}
	applyTheChanges := false // Doesn't matter in this case as no changes will be applied anyway

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges, &stdout, &stderr)

	// Check results
	expectedStdout := "" // No differences should be reported

	assert.Equal(t, expectedStdout, stdout.String())
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestMysqlConfDiff_NoDiff_ON_1(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()

	confOptions := map[string]any{"key1": "1"}
	serverVariables := map[string]any{"key1": "ON"}
	applyTheChanges := false // Doesn't matter in this case as no changes will be applied anyway

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges, &stdout, &stderr)

	// Check results
	expectedStdout := "" // No differences should be reported

	assert.Equal(t, expectedStdout, stdout.String())
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestMysqlConfDiff_NoDiff_OFF_0(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()

	confOptions := map[string]any{"key1": "0"}
	serverVariables := map[string]any{"key1": "OFF"}
	applyTheChanges := false // Doesn't matter in this case as no changes will be applied anyway

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges, &stdout, &stderr)

	// Check results
	expectedStdout := "" // No differences should be reported

	assert.Equal(t, expectedStdout, stdout.String())
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestMysqlConfDiff_NoDiff_DirectorySlash(t *testing.T) {
	// Prepare dependencies and inputs
	conn, m, err := sqlmock.New()
	require.NoError(t, err)
	db := &dbConn{conn: conn}
	defer db.close()

	confOptions := map[string]any{"key1": "/my/dir"}
	serverVariables := map[string]any{"key1": "/my/dir/"}
	applyTheChanges := false // Doesn't matter in this case as no changes will be applied anyway

	// Capture stdout and stderr
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// Run function
	mysqlConfDiff(db, confOptions, serverVariables, applyTheChanges, &stdout, &stderr)

	// Check results
	expectedStdout := "" // No differences should be reported

	assert.Equal(t, expectedStdout, stdout.String())
	assert.Empty(t, stderr.String())
	require.NoError(t, m.ExpectationsWereMet())
}

func TestNormalizeKeys(t *testing.T) {
	input := map[string]any{
		"key-test1": "value1",
		"key-test2": "value2",
		"keyTest3":  "value3",
		"KEYTEST4":  "value4",
	}
	expected := map[string]any{
		"KEY_TEST1": "value1",
		"KEY_TEST2": "value2",
		"KEYTEST3":  "value3",
		"KEYTEST4":  "value4",
	}

	result := normalizeKeys(input)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, but got %v", expected, result)
	}
}

func TestLimitToWatchedOptions(t *testing.T) {
	fullOptions := map[string]any{
		"key1": "1",
		"key2": "2",
		"key3": "3",
	}
	watchedOptions := map[string]any{
		"key1": true,
		"key3": true,
	}
	expected := map[string]any{
		"KEY1": "1",
		"KEY3": "3",
	}

	result := limitToWatchedOptions(fullOptions, watchedOptions)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("limitToWatchedOptions() = %v, want %v", result, expected)
	}
}
