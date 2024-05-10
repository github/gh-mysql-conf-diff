package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestConvertOptionNameToVariableKey(t *testing.T) {
	input := "my-option-name"
	expected := "MY_OPTION_NAME"
	result := GetVariableKeyFrom(input)
	if result != expected {
		t.Fatalf("got %s, want %s", result, expected)
	}
}

func TestApplySetting_Int_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	d := &dbConn{conn: db}

	mock.ExpectExec(`SET GLOBAL MAX_CONNECTIONS = \?`).
		WithArgs(1000).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = d.applySetting("MAX_CONNECTIONS", "1000")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestApplySetting_Str_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	d := &dbConn{conn: db}

	mock.ExpectExec(`SET GLOBAL CHARACTER_SET_SERVER = \?`).
		WithArgs("utf8mb4").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = d.applySetting("CHARACTER_SET_SERVER", "utf8mb4")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
