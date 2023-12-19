// Utility: db-mysql-conf-diff
//
// This utility prints a diff between `my.cnf` on disk with a running MySQL
// server and optionally applies the changes found to the server. The program
// honors version blocks (e.g. `[mysql-8.0]` blocks will only be checked for
// servers running `8.0.*`).
//
// A simple run of the utility might be as follows:
//
//	$ db-mysql-conf-diff /etc/mysql/my.cnf localhost:3306
//
//	Difference found for: CONNECT_TIMEOUT
//	  my.cnf:    60
//	  mysqld:    30
//
// By default the utility runs in read only (informational mode). To apply the
// changes, use the `--apply-changes` flag. This is not enabled by default. If
// you run `--apply-changes` you need to use `--watch-options“ as well:
//
//	$ db-mysql-conf-diff /etc/mysql/my.cnf localhost:3306 \
//	   --watch-options connect_timeout,delay_key_write --apply-changes
//
// The program needs to connect to MySQL with a user that has the correct
// permissions. The username and password combo can be set using environment
// variables `$MYSQL_USER` and `$MYSQL_PASSWORD`.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	os.Exit(runWithReturnCode())
}

func runWithReturnCode() int {
	// Handle program configuration setup by parsing arguments.
	cli := newInputContext()
	context, err := cli.parseArgs(os.Args[1:])
	if err != nil {
		if errors.Is(err, errHelpFlagIsSet) {
			_, _ = fmt.Fprintf(os.Stderr, "%s", cli.getHelpMessage())
			return 0
		}
		_, _ = fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
		return 1
	}
	// If --apply-changes, then fail if no --watch-options
	if context.applyTheChanges && len(context.optionKeysToWatch) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Fatal: --watch-options is required when using --apply-changes\n")
		return 1
	}
	// Get the DB connection
	db, err := getDB(context)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}
	defer db.close()
	// Get the two option maps, one from my.cnf, and one from the
	// server variables for comparison.
	confOptions, serverVariables, err := getOptionsFrom(context.configPath, db)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}
	// If --watch-options is set with any values, use it as a filter to limit
	// the options to these values. If it is not set, then all options are
	// used.
	if len(context.optionKeysToWatch) > 0 {
		// Limit the my.cnf options to those specified by --watch-options
		confOptions = limitToWatchedOptions(confOptions, context.optionKeysToWatch)
	} else {
		// Limit the my.cnf options to the available server variables.
		// If we don't limit to this, it will generate a lot of warnings,
		// because the my.cnf file allows additional options (e.g. `USER` or
		// `REPLICATE_SAME_SERVER_ID`) than the server variables.
		confOptions = limitToWatchedOptions(confOptions, serverVariables)
	}
	// Compare the options maps and print results to stdout and stderr
	// as appropriate. If --apply-changes, then also apply the changes
	// to the server.
	mysqlConfDiff(
		db, confOptions, serverVariables, context.applyTheChanges,
		os.Stdout, os.Stderr)
	return 0
}

// Given the connection information defined in the run context, this
// function connects to the MySQL server and returns an open connection.
func getDB(context *RunContext) (db *dbConn, err error) {
	// Get user and password information
	user, password, err := getMySQLUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get MySQL user info: %w", err)
	}

	// Connect to the MySQL server
	db, err = connect(fmt.Sprintf(
		"%s:%s@tcp(%s)/", user, password, context.serverAndPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	return db, nil
}

// Given the my.cnf path and a database connection, this function reads
// the my.cnf file and queries the server for its variables. It returns
// these as two maps.
func getOptionsFrom(configPath string, db *dbConn) (
	confOptions map[string]any, serverVariables map[string]any, err error) {
	// Get the running MySQL version. This is necessary to interpret
	// the configuration option blocks correctly.
	version, err := db.getVersion()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to read mysql version: %w", err)
	}
	// Get the variables of the running server.
	serverVariables, err = db.getVariables()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query MySQL for server variables: %w", err)
	}
	// Read my.cnf configuration file
	mysqlConfig, err := NewMySQLConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load MySQL config: %w", err)
	}
	// Limit the my.cnf options to those for the running MySQL version
	confOptions = mysqlConfig.ComposeForVersion(version)
	return confOptions, serverVariables, nil
}

// Given the my.cnf options map and server variables map, this function
// compares the two and prints any differences to stdout. If the
// --apply-changes flag is set, then the function will also apply the
// changes to the server and print the change it made.
func mysqlConfDiff(
	db *dbConn,
	confOptions map[string]any,
	serverVariables map[string]any,
	applyTheChanges bool,
	stdout, stderr io.Writer,
) {
	// Loop through the my.cnf options and compare to the server variables
	// watched by the user.
	for key, optionValue := range confOptions {
		// If the option is not in the server variables, then it is
		// potentially invalid. Report this to the user.
		serverValue, keyExists := serverVariables[key].(string)
		if !keyExists {
			_, _ = fmt.Fprintf(stderr,
				"Warning: option '%s' in configuration file was "+
					"not found in server variables\n", key)
			continue
		}
		if serverValue == optionValue {
			continue // Nothing to do
		}
		// Handle ON and OFF in server variables' equality with 1
		// and 0, respectively.
		if serverValue == "ON" && optionValue == "1" {
			continue // Nothing to do
		}
		if serverValue == "OFF" && optionValue == "0" {
			continue // Nothing to do
		}
		// Handle directories that end with a slash
		if strings.HasSuffix(serverValue, "/") && serverValue[:len(serverValue)-1] == optionValue {
			continue // Nothing to do
		}
		// Report on any differences to console user
		_, _ = fmt.Fprintf(stdout, "Difference found for: %s\n", key)
		_, _ = fmt.Fprintf(stdout, "  my.cnf:    %s\n", optionValue)
		_, _ = fmt.Fprintf(stdout, "  mysqld:    %s\n", serverValue)
		// If the --apply-changes flag is provided, actually apply the changes
		if applyTheChanges {
			err := db.applySetting(
				key, optionValue)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "Warning: Failed to SET variable: %v\n", err)
			}
			_, _ = fmt.Fprintf(stdout, "Set variable:\n  %s = %s\n", key, optionValue)
		}
	}
}

// Given a map of options, this function returns a new map with the
// keys normalized to the format used by the MySQL server.
func normalizeKeys(input map[string]any) map[string]any {
	result := make(map[string]any)
	for key := range input {
		normalizedKey := GetVariableKeyFrom(key)
		result[normalizedKey] = input[key]
	}
	return result
}

// Only watch certain settings, based on --watch-options.
// This function effectively limits the original map to only the keys
// that are in the watchedOptions map.
func limitToWatchedOptions(
	fullOptions map[string]any,
	watchedOptions map[string]any,
) map[string]any {
	fullOptions = normalizeKeys(fullOptions)
	watchedOptions = normalizeKeys(watchedOptions)

	limitedOptions := make(map[string]any)
	for key := range fullOptions {
		if _, ok := watchedOptions[key]; ok {
			limitedOptions[key] = fullOptions[key]
		}
	}
	return limitedOptions
}
