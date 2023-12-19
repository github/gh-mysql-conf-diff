## Contributing

[fork]: https://github.com/github/mysql-conf-diff/fork
[pr]: https://github.com/github/mysql-conf-diff/compare
[style]: https://github.com/github/mysql-conf-diff/blob/main/.golangci.yaml
[code-of-conduct]: CODE_OF_CONDUCT.md

Hi there! We're happy that you'd like to contribute to this project. Your help is essential for keeping it great.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE.md).

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

**Environment Setup and External Dependencies**

_Environment Requirements_:
To run mysql-conf-diff, you need an environment capable of running Go (Golang) applications. The minimum requirements are:

1. **Go (Golang) Environment**: The tool is developed in Go, so you need to have Go installed on your system. Ensure you have the latest stable version of Go for optimal performance and compatibility.
2. **MySQL Server Access**: As the tool interacts with MySQL servers, you must have network access to the MySQL server you wish to compare configurations against.

_Setting Up the Environment_:
1. **Install Go**: Download and install the Go language runtime from the [official Go website](https://golang.org/dl/). Follow the installation instructions specific to your operating system.
2. **Set Go Environment Variables**: Configure your Go workspace by setting the `GOPATH` and `GOROOT` environment variables as per the Go documentation.
3. **Clone the Repository**: Clone the mysql-conf-diff repository from its source to your local Go workspace.
4. **Build the Tool**: Navigate to the cloned repository directory and run `go build` to compile the application.

_External Dependencies_:
mysql-conf-diff has several external dependencies which need to be installed:

1. **MySQL Client Libraries**: The tool requires MySQL client libraries for database communication. Install these libraries based on your operating system's package manager.
2. **Go MySQL Driver**: This is a Go-based MySQL driver needed for database interactions. It can be installed using Go's package manager with `go get -u github.com/go-sql-driver/mysql`.

After setting up the environment and installing the necessary dependencies, you should be able to run mysql-conf-diff successfully to compare and synchronize MySQL configurations.

Be sure to follow the [GitHub logo guidelines](https://github.com/logos).


## Prerequisites for running and testing code

These are one time installations required to be able to test your changes locally as part of the pull request (PR) submission process.

1. install Go [through download](https://go.dev/doc/install) | [through Homebrew](https://formulae.brew.sh/formula/go)
1. [install golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Building and testing code

### Compiling the Binary

To compile the binary from the source, from the root of the project run the following command. This will create an executable in the `bin/mysql-conf-diff/` directory.

   ```sh
   go build -o bin/mysql-conf-diff/mysql-conf-diff ./cmd/mysql-conf-diff/
   ```

### Running unit tests

To run the unit test suite:

   ```sh
   go test ./...
   ```

## Coding Conventions

Our codebase adheres to certain coding conventions. Before you contribute, please make sure to follow them:

1. **Formatting**: We use `gofmt` to automatically format our code. Please make sure to run `gofmt -s -w .` on your code before committing.

2. **Naming**: We prefer short, concise names for local variables and more descriptive names for exported functions and variables. Acronyms should be all uppercase.

3. **Error Handling**: Always check errors and handle them immediately. Do not ignore errors or use panic for normal error handling.

4. **Comments**: Write a comment for every exported function, variable, and type. The comment should start with the name of the thing it's describing.

5. **Packages**: Each package should have a single purpose and provide a clean, simple API. The package name should be a noun, and the functions in the package should be actions on that noun.

6. **Testing**: Write tests for your code. Test functions should be named `TestXxx`, where `Xxx` is the name of the function being tested.

We use several tools to maintain the quality of our codebase:

- `golint`: Checks the code for style issues. Run it on your code to make sure it adheres to our style guide.
- `go vet`: Checks the code for common errors. Run it on your code to catch any potential issues.
- `staticcheck`: A static analysis tool that checks for a wide range of issues. We recommend running it on your code.
- `golangci-lint`: A fast Go linters runner. We use it in our CI/CD pipeline to catch any issues before they get merged.

Remember, these are conventions and tools to help write better code. Always use your best judgment and consider the specific needs of your project.

## Submitting a pull request

1. [Fork][fork] and clone the repository
1. Configure and install the dependencies: `script/bootstrap`
1. Make sure the tests pass on your machine: `go test -v ./...`
1. Make sure linter passes on your machine: `golangci-lint run`
1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and make sure the tests and linter still pass
1. Push to your fork and [submit a pull request][pr]
1. Pat yourself on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Follow the [style guide][style].
- Write tests.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)