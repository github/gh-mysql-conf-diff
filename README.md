## mysql-conf-diff

_Description_: 

`mysql-conf-diff` is a utility tool designed for database administrators and developers working with MySQL. It compares my.cnf configuration files on disk with settings of a running MySQL server, highlighting any differences. It's particularly useful for ensuring consistency in MySQL server configurations. Written in Go, it differs from similar tools by its ability to apply detected configuration changes directly to the server and its support for version-specific configuration blocks.

_Current project status_:

The project is open to pull requests and the utility is actively used in production at GitHub.

_Features_:
1. Compares `my.cnf` files with live MySQL server settings.
2. Supports version-specific configuration blocks (e.g., `[mysql-8.0]`).
3. Read-only informational mode by default.
4. Option to apply changes to the server using `--apply-changes` flag.
5. Allows specification of which options to watch and apply using `--watch-options`.
6. User authentication through environment variables `$MYSQL_USER` and `$MYSQL_PASSWORD`.

_Limitations_:
1. Requires user authentication with appropriate permissions.
2. Only compatible with MySQL servers.
3. Manual specification of options to watch when applying changes.

_Goals and Scope_:
The primary goal of `mysql-conf-diff` is to provide a reliable and efficient way to ensure consistency between MySQL server configurations and `my.cnf` files. It aims to streamline the configuration management process, reduce errors, and save time for database administrators. The utility is scoped to focus on comparing and optionally synchronizing configurations, without delving into other aspects of database management.

## Example Usage

By default the utility runs in read only (informational mode). To apply the
changes, use the `--apply-changes` flag. This is not enabled by default. If
you run `--apply-changes` you need to use `--watch-options` as well:

	$ db-mysql-conf-diff /etc/mysql/my.cnf localhost:3306 \
	   --watch-options connect_timeout,delay_key_write --apply-changes

## Background 

**Development Roadmap and Contributions**

_Development Roadmap_:
For a detailed view of our open issues, please refer to our [tracker](https://github.com/github/mysql-conf-diff/issues). Please feel free to submit feature requests or bugs.

_Contributions_:
We highly value contributions from the community and encourage developers, database administrators, and other interested individuals to contribute to `mysql-conf-diff`. Whether it's by reporting bugs, suggesting enhancements, or submitting code changes, your input is important to the growth and improvement of this tool. 

For detailed guidelines on how to contribute, please see our [CONTRIBUTING.md](CONTRIBUTING.md). This document provides all the information you need to get started with contributing to `mysql-conf-diff`, including coding standards, pull request processes, and how to set up your development environment.

Your contributions are welcomed and greatly appreciated.

## Requirements

**Environment Setup and External Dependencies**

_Environment Requirements_:
To run db-mysql-conf-diff, you need an environment capable of running Go (Golang) applications. The minimum requirements are:

1. **Go (Golang) Environment**: The tool is developed in Go, so you need to have Go installed on your system. Ensure you have the latest stable version of Go for optimal performance and compatibility.
2. **MySQL Server Access**: As the tool interacts with MySQL servers, you must have network access to the MySQL server you wish to compare configurations against.

_Setting Up the Environment_:
1. **Install Go**: Download and install the Go language runtime from the [official Go website](https://golang.org/dl/). Follow the installation instructions specific to your operating system.
2. **Set Go Environment Variables**: Configure your Go workspace by setting the `GOPATH` and `GOROOT` environment variables as per the Go documentation.
3. **Clone the Repository**: Clone the db-mysql-conf-diff repository from its source to your local Go workspace.
4. **Build the Tool**: Navigate to the cloned repository directory and run `go build` to compile the application.

_External Dependencies_:
db-mysql-conf-diff has several external dependencies which need to be installed:

1. **MySQL Client Libraries**: The tool requires MySQL client libraries for database communication. Install these libraries based on your operating system's package manager.
2. **Go MySQL Driver**: This is a Go-based MySQL driver needed for database interactions. It can be installed using Go's package manager with `go get -u github.com/go-sql-driver/mysql`.

After setting up the environment and installing the necessary dependencies, you should be able to run db-mysql-conf-diff successfully to compare and synchronize MySQL configurations.

Be sure to follow the [GitHub logo guidelines](https://github.com/logos).

## License 

This project is licensed under the terms of the MIT open source license. Please refer to [MIT](./LICENSE.txt) for the full terms.

## Maintainers 

[@adamsc64](https://github.com/adamsc64)

## Support

## Support

Here are our support expectations:

1. **Community Support**: Our primary support channel is through our community. Users are encouraged to share their experiences, solutions, and best practices among each other. For general questions or discussions, please use [community forum link or communication channel].
1. **Issue Tracker**: For reporting bugs or requesting new features, please use our issue tracker [link to issue tracker]. This platform is regularly monitored by our team, and we aim to respond to each issue appropriately.
1. **Contribution**: Users who wish to contribute fixes or improvements are welcome to do so. Please refer to our contribution guidelines [link to CONTRIBUTING.md] for more information on how to contribute.

Please note that as an open-source project, support is largely dependent on the availability and capacity of our community and development team. We appreciate your understanding and patience.

## Acknowledgement

This acknowledgment is a token of appreciation for the support, guidance, and resources provided by GitHub and its Database Infrastructure Team, which have been pivotal in the development and continued improvement of `mysql-conf-diff`.
