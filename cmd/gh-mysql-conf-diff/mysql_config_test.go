package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMySQLConfig(t *testing.T) {
	data := []byte(`
[section1]
key1=value1

[section2]
key2=value2
`)
	config, err := NewMySQLConfig(data)
	require.NoError(t, err)

	if len(config.sectionTitles) != 3 ||
		config.sectionTitles[0] != "DEFAULT" ||
		config.sectionTitles[1] != "section1" ||
		config.sectionTitles[2] != "section2" {
		t.Fatalf("Expected section titles unexpected, got: %v", config.sectionTitles)
	}
}

func TestComposeForVersion57(t *testing.T) {
	cfg, err := NewMySQLConfig(
		[]byte(`
[mysqld]
key1=value1

[mysqld-5.7]
key1=value57

[mysqld-8.0]
key1=value80
`),
	)
	require.NoError(t, err)

	version := MySQLVersion{
		Major: 5,
		Minor: 7,
		Patch: 34,
	}

	expected := map[string]any{
		"key1": "value57",
	}

	actual := cfg.ComposeForVersion(version)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("ComposeForVersion() = %v, want %v", actual, expected)
	}
}

func TestComposeForVersion80(t *testing.T) {
	cfg, err := NewMySQLConfig(
		[]byte(`
[mysqld]
key1=value1

[mysqld-5.7]
key1=value57

[mysqld-8.0]
key1=value80
`),
	)
	require.NoError(t, err)

	version := MySQLVersion{
		Major: 8,
		Minor: 0,
		Patch: 28,
	}

	expected := map[string]any{
		"key1": "value80",
	}

	actual := cfg.ComposeForVersion(version)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("ComposeForVersion() = %v, want %v", actual, expected)
	}
}

func TestNoCompose(t *testing.T) {
	cfg, err := NewMySQLConfig(
		[]byte(`
[mysqld]
key1=value1

[mysqld-5.7]
key2=value57

[mysqld-8.0]
key2=value80
`),
	)
	require.NoError(t, err)

	version := MySQLVersion{
		Major: 5,
		Minor: 7,
		Patch: 34,
	}

	actual := cfg.ComposeForVersion(version)
	if value, ok := actual["key1"]; ok {
		if value != "value1" {
			t.Fatalf("Unexpected value for key1, got: %v, want: %v", value, "value1")
		}
	} else {
		t.Fatalf("key1 not found in map")
	}
}

func TestVersionMatchWithCatchall(t *testing.T) {
	version := MySQLVersion{
		Major: 5,
		Minor: 7,
		Patch: 34,
	}
	sectionTitle := "mysqld"
	expected := true
	result := isOptionBlockMatch(version, sectionTitle)
	require.Equal(t, expected, result)
}

func TestVersionMatch(t *testing.T) {
	version := MySQLVersion{
		Major: 5,
		Minor: 7,
		Patch: 34,
	}
	sectionTitle := "mysqld-5.7"
	expected := true
	result := isOptionBlockMatch(version, sectionTitle)
	require.Equal(t, expected, result)
}

func TestVersionNotMatch(t *testing.T) {
	version := MySQLVersion{
		Major: 5,
		Minor: 7,
		Patch: 34,
	}
	sectionTitle := "mysqld-8.0"
	expected := false
	result := isOptionBlockMatch(version, sectionTitle)
	require.Equal(t, expected, result)
}

func TestCleanRequireEquals(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "when line does not contain '='",
			input:    []byte("[test]\nfoo\n"),
			expected: []byte("[test]\n"),
		},
		{
			name:     "when line contains '='",
			input:    []byte("[test]\nfoo = bar\n"),
			expected: []byte("[test]\nfoo = bar\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clean(tt.input)
			require.Equal(t, string(tt.expected), string(result))
		})
	}
}

func TestCleanIncludes(t *testing.T) {
	input := []byte("[test]\nfoo = bar\n!includedir /etc/mysql/conf.d/\n")
	expected := []byte("[test]\nfoo = bar\n")

	result := clean(input)
	require.Equal(t, string(expected), string(result))
}

func TestNormalizeSizes(t *testing.T) {
	require.Equal(t, "1024", normalize("1K"))
	require.Equal(t, "1048576", normalize("1M"))
	require.Equal(t, "1073741824", normalize("1G"))
	require.Equal(t, "1099511627776", normalize("1T"))

	require.Equal(t, "5120", normalize("5K"))
	require.Equal(t, "5242880", normalize("5M"))
	require.Equal(t, "5368709120", normalize("5G"))
	require.Equal(t, "5497558138880", normalize("5T"))
}

func TestNormalizeNoChange(t *testing.T) {
	require.Equal(t, "1GK", normalize("1GK"))
	require.Equal(t, "K", normalize("K"))
	require.Equal(t, "1024", normalize("1024"))
	require.Equal(t, "string-value", normalize("string-value"))
	require.Equal(t, "string-valueG", normalize("string-valueG"))
}
