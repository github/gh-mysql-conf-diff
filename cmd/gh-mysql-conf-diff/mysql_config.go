package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

// MySQLConfig represents a loaded MySQL config file.
type MySQLConfig struct {
	sectionTitles []string
	cfg           *ini.File
}

// NewMySQLConfig creates a new MySQLConfig object from the given config,
// which can be either the path to the config file or []byte with the config
// file contents.
func NewMySQLConfig(config any) (*MySQLConfig, error) {
	// Read the config file contents and handle polymorphic type
	confContents, err := readFile(config)
	if err != nil {
		return nil, err
	}
	// Remove unsupported lines and directives
	confContents = clean(confContents)
	// Parse the resulting, cleaned config
	cfg, err := ini.Load(confContents)
	if err != nil {
		return nil, err
	}
	// Get the section titles
	sectionTitles := cfg.SectionStrings()
	// Construct and return the new object
	return &MySQLConfig{
		sectionTitles: sectionTitles,
		cfg:           cfg,
	}, nil
}

// ComposeForVersion composes a map of all the MySQL config settings that
// should be applied for the given MySQL version. The map keys are the MySQL
// config option names.
func (c *MySQLConfig) ComposeForVersion(version MySQLVersion) map[string]any {
	allSettings := make(map[string]any)
	for _, sectionTitle := range c.sectionTitles {
		if !isOptionBlockMatch(version, sectionTitle) {
			continue
		}
		settings := c.cfg.Section(sectionTitle).KeysHash()

		for key, value := range settings {
			allSettings[key] = normalize(value)
		}
	}
	return allSettings
}

func isOptionBlockMatch(version MySQLVersion, sectionTitle string) bool {
	if sectionTitle == "mysqld" {
		return true
	}
	if sectionTitle == fmt.Sprintf("mysqld-%d.%d", version.Major, version.Minor) {
		return true
	}
	return false
}

// Reads the config file contents from the given config, which can be
// either the path to the config file or []byte with the config file
// contents.
func readFile(config any) ([]byte, error) {
	var confContents []byte
	var err error
	switch c := config.(type) {
	case []byte:
		confContents = c
	case string:
		// Read the config file contents
		confContents, err = os.ReadFile(c)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type given to NewMySQLConfig")
	}
	return confContents, nil
}

// Removes lines from the config file that are not supported by the
// utility.
func clean(configContents []byte) []byte {
	lines := strings.Split(string(configContents), "\n")
	var newLines []string
	for _, line := range lines {
		isSection := strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
		isEmpty := strings.TrimSpace(line) == ""
		hasEquals := strings.Contains(line, "=")
		isIncludeDirective := strings.HasPrefix(line, "!include")
		// Remove options that don't have a '=' in them (boolean options)
		// This is a workaround for the fact that the go-ini library
		// handles these options ambiguously.
		if !isSection && !isEmpty && !hasEquals {
			continue
		}
		// This utility does not support include directives at the
		// moment, so remove them.
		if isIncludeDirective {
			continue
		}
		newLines = append(newLines, line)
	}
	return []byte(strings.Join(newLines, "\n"))
}

// Normalizes values that are allowed in my.cnf but are not allowed in
// SET GLOBAL statements.
func normalize(value string) string {
	if value == "" {
		return value
	}
	sizeSuffixes := []string{"K", "M", "G", "T"}
	for _, suffix := range sizeSuffixes {
		if strings.HasSuffix(value, suffix) {
			newValue, err := normalizeDataSize(value)
			if err != nil {
				return value
			}
			return newValue
		}
	}
	return value
}

// Normalizes data size values, e.g. turns "1K" into "1024".
func normalizeDataSize(value string) (string, error) {
	sizes := map[string]int64{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": 1024 * 1024 * 1024 * 1024,
	}
	// Split the value into magnitude and suffix
	magnitudeStr := value[:len(value)-1]
	suffix := value[len(value)-1:]
	// Convert the magnitude to an integer for multiplication
	magnitude, err := strconv.ParseInt(magnitudeStr, 10, 64)
	if err != nil {
		return "", err
	}
	// Multiply the magnitude by the suffix multiplier
	suffix = strings.ToUpper(suffix)
	if multiplier, ok := sizes[suffix]; ok {
		bytesValue := magnitude * multiplier
		return strconv.FormatInt(bytesValue, 10), nil
	}
	// If the suffix is not valid, return an error
	return "", errors.New("invalid data size suffix, expected 'K', 'M', 'G', or 'T'")
}
