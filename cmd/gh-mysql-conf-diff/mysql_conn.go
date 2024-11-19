package main

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var keyValidator = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type dbConn struct {
	conn    *sql.DB
	timeout time.Duration
}

func connect(user string, password string, host string, port int, timeout time.Duration) (*dbConn, error) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?timeout=%ds", user, password, host, port, timeout/time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &dbConn{conn: db, timeout: timeout}, nil
}

func (db *dbConn) getOne(query string, args ...any) (answer string, err error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel() // should be called even on success

	var rows *sql.Rows
	rows, err = db.conn.QueryContext(timeoutContext, query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&answer)
	if err != nil {
		return "", err
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (db *dbConn) getKeyValues(query string, args ...any) (map[string]any, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel() // should be called even on success

	var rows *sql.Rows
	rows, err := db.conn.QueryContext(timeoutContext, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the MySQL configuration variables into a map
	keyValues := make(map[string]any)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		keyValues[key] = value
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return keyValues, nil
}

func (db *dbConn) exec(query string, args ...any) (sql.Result, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel() // should be called even on success

	return db.conn.ExecContext(timeoutContext, query, args...)
}

func (db *dbConn) close() error {
	return db.conn.Close()
}

// Gets MySQL version from server and returns it as a rich object.
func (db *dbConn) getVersion() (MySQLVersion, error) {
	firstResult, err := db.getOne("SELECT VERSION()")
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
	serverVariables, err := db.getKeyValues("SHOW VARIABLES")
	if err != nil {
		return nil, err
	}
	// normalize all the keys as uppercase
	for key, value := range serverVariables {
		serverVariables[strings.ToUpper(key)] = value
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
		_, err := db.exec(fmt.Sprintf(`SET GLOBAL %s = ?`, key), valueInt)
		if err != nil {
			return err
		}
	} else {
		// treating as a string
		_, err := db.exec(fmt.Sprintf(`SET GLOBAL %s = ?`, key), valueStr)
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
