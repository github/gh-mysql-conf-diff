package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

// RunContext contains the information needed to run the program.
type RunContext struct {
	configPath        string
	serverAndPort     string
	optionKeysToWatch map[string]any
	applyTheChanges   bool
}

// InputContext contains the information from the command-line arguments.
type InputContext struct {
	optionsToWatchFlag []string
	executeFlag        bool
	helpFlag           bool

	positionals []string

	flagset *pflag.FlagSet
}

var errHelpFlagIsSet = errors.New("help flag is set")

func newInputContext() *InputContext {
	cli := InputContext{flagset: pflag.NewFlagSet("", pflag.ContinueOnError)}
	cli.flagset.StringSliceVarP(&cli.optionsToWatchFlag, "watch-options", "", nil,
		"A comma-separated list of MySQL config file option names to watch")
	cli.flagset.BoolVarP(&cli.executeFlag, "apply-changes", "", false,
		"If provided, actually apply the changes discovered. [optional]")
	cli.flagset.BoolVarP(&cli.helpFlag, "help", "h", false, "Print this help message and exit")
	cli.flagset.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, cli.getHelpMessage())
	}
	return &cli
}

func (c *InputContext) parseArgs(args []string) (
	context *RunContext, err error) {
	// Parse the command-line arguments
	err = c.flagset.Parse(args)
	if err != nil {
		return nil, err
	}
	c.positionals = c.flagset.Args()
	// Validate the command-line arguments
	if c.helpFlag {
		return nil, errHelpFlagIsSet
	}
	if len(c.positionals) != 2 {
		return nil, fmt.Errorf("invalid number of positional arguments")
	}
	if c.executeFlag && len(c.optionsToWatchFlag) == 0 {
		return nil, fmt.Errorf("--watch-options required when running --apply-changes")
	}
	// Handle the --watch-options flag
	optionsToWatch := make(map[string]any)
	for _, option := range c.optionsToWatchFlag {
		optionsToWatch[option] = true
	}
	return &RunContext{
		configPath:        c.positionals[0],
		serverAndPort:     c.positionals[1],
		optionKeysToWatch: optionsToWatch,
		applyTheChanges:   c.executeFlag,
	}, nil
}

// Returns the help message to display to the user.
func (c *InputContext) getHelpMessage() string {
	var message strings.Builder

	_, _ = fmt.Fprint(&message, "Usage: ", getBinaryName(), " <path_to_my.cnf> <server:port> "+
		"[--watch-options option1,option2,option3 [--apply-changes]]")
	_, _ = fmt.Fprint(&message, "\n\n")
	_, _ = fmt.Fprint(&message,
		"This utility checks the MySQL configuration on disk against the server variable "+
			"settings of a running MySQL server and prints a diff summary to stdout. It only "+
			"inspects the specified MySQL configuration options. Version option blocks (i.e. "+
			"[mysqld-5.7] and [mysqld-8.0]) are honored given the MySQL server version of the "+
			"server provided. The program can optionally *apply* changes found onto the "+
			"running MySQL server."+
			"\n\n"+
			"Set environment variable $MYSQL_USER and $MYSQL_PASSWORD to specify connection "+
			"information.")
	_, _ = fmt.Fprint(&message, "\n\n")
	_, _ = fmt.Fprint(&message, c.flagset.FlagUsages())

	return message.String()
}

// Returns the username and password to use to connect to the MySQL server.
func getMySQLUserInfo() (username, password string, err error) {
	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	if username == "" {
		return "", "", fmt.Errorf("no user provided")
	}
	if password == "" {
		return "", "", fmt.Errorf("no password provided")
	}

	return username, password, nil
}

// Returns the name of the binary, without the path.
func getBinaryName() string {
	binaryName := os.Args[0]
	binaryName = filepath.Base(binaryName)

	return binaryName
}
