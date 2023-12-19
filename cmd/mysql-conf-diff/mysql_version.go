package main

import (
	"fmt"
	"strconv"
	"strings"
)

// MySQLVersion is a rich object representing a MySQL version.
type MySQLVersion struct {
	Major int
	Minor int
	Patch int
}

// Given a version string, parse it into a MySQLVersion object.
func ParseVersion(version string) (MySQLVersion, error) {
	// Version numbers can have a suffix like `5.7.34-log`, i.e. 5.7.34 with
	// extra logging enabled. There are other possible suffixes to the
	// version number besides -log, including -opt, -rc, -beta, etc.
	// For the purposes of this tool, we don't care about those, so ignore.
	if strings.Contains(version, "-") {
		version = strings.Split(version, "-")[0]
	}
	// Split the version into parts and validate the format.
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return MySQLVersion{}, fmt.Errorf("invalid version format")
	}
	// Convert the parts into integers.
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return MySQLVersion{}, fmt.Errorf("invalid major version: %w", err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return MySQLVersion{}, fmt.Errorf("invalid minor version: %w", err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return MySQLVersion{}, fmt.Errorf("invalid patch version: %w", err)
	}
	// Return the composed final object.
	return MySQLVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}
