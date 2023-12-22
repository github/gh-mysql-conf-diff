package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var keyValidator = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type dbConn struct {
	conn *sql.DB
}

func connect(dataSourceName string) (*dbConn, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &dbConn{conn: db}, nil
}

func (db *dbConn) close() error {
	return db.conn.Close()
}

// Gets MySQL version from server and returns it as a rich object.
func (db *dbConn) getVersion() (MySQLVersion, error) {
	rows, err := db.conn.Query("SELECT VERSION()")
	if err != nil {
		return MySQLVersion{}, err
	}
	defer rows.Close()
	rows.Next()
	var firstResult string
	err = rows.Scan(&firstResult)
	if err != nil {
		return MySQLVersion{}, err
	}
	err = rows.Err()
	if err != nil {
		return MySQLVersion{}, err
	}
	version, err := ParseVersion(firstResult)
	if err != nil {
		return MySQLVersion{}, err
	}

	return version, nil
}

// Get MySQL configuration variables.
func (db *dbConn) getVariables() (map[string]any, error) {
	//nolint:execinquery // SHOW is incorrectly failing the lint
	rows, err := db.conn.Query("SHOW VARIABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the MySQL configuration variables into a map
	serverVariables := make(map[string]any)
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		serverVariables[strings.ToUpper(key)] = value
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return serverVariables, nil
}

// Apply a change of a setting to the MySQL server.
func (db *dbConn) applySetting(key string, value any) error {
	// ensure that submitted data only contains certain subset of symbols
	if !keyValidator.MatchString(key) {
		return fmt.Errorf("invalid key: %s", key)
	}
	// initially, we assume that the value is a string
	valueStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type: %T", value)
	}
	if valueInt, err := strconv.Atoi(valueStr); err == nil {
		// converted to an int successfully, so we treat it as an int
		_, err := db.conn.Exec(fmt.Sprintf(`SET GLOBAL %s = ?`, key), valueInt)
		if err != nil {
			return err
		}
	} else {
		// treating as a string
		_, err := db.conn.Exec(fmt.Sprintf(`SET GLOBAL %s = ?`, key), valueStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetVariableKeyFrom converts the key name from mysql configuration
// format to match the MySQL server variable key format.
func GetVariableKeyFrom(optionName string) string {
	normalizedKey := strings.ToUpper(optionName)
	normalizedKey = strings.ReplaceAll(normalizedKey, "-", "_")
	return normalizedKey
}
