# Contributing

[fork]: https://github.com/github/gh-mysql-conf-diff/fork
[pr]: https://github.com/github/gh-mysql-conf-diff/compare
[style]: https://github.com/github/gh-mysql-conf-diff/blob/main/.golangci.yml

Hi there! We're happy that you'd like to contribute to this project.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE.txt). Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Environment Requirements
To run `gh-mysql-conf-diff`, you need an environment capable of running Go (Golang) applications. The minimum requirements are:

1. **Go (Golang) Environment**: The tool is developed in Go, so you need to have Go installed on your system. Ensure you have at least the version of Go in [go.mod](go.mod) for optimal compatibility.
2. **MySQL Server Access**: As the tool interacts with MySQL servers, you must have network access to a MySQL server you wish to compare configurations against. This can be a MySQL server running on `localhost`.

### Setting Up the Environment

1. **Install Prerequisites**: Follow the instructions in the [README](README.md#requirements) how to install Go and mysql client libraries.
1. **Clone the Repository**: Instead of running `go install`, clone the `gh-mysql-conf-diff` repository to your local development environment.

After setting up the environment and installing the necessary dependencies, you should be able to build `gh-mysql-conf-diff` successfully.

## Building and testing code

### Running unit tests

Before building the binary, you might want to run the unit test suite to ensure everything is set up correctly with your Go environment. To do this:

   ```sh
   go test ./...
   ```

### Compiling the Binary

To compile the binary from the source, from the root of the project run the following command. This will create an executable in the `bin/gh-mysql-conf-diff/` directory.

   ```sh
   go build -o bin/gh-mysql-conf-diff/gh-mysql-conf-diff ./cmd/gh-mysql-conf-diff/
   ```

### Starting a development database with docker-compose

The repository contains helpful a helpful docker-compose configuration to start up mysqld in a container for testing and development work. To start mysql, run this command from the project root:

   ```sh
   docker-compose up -d mysql-database
   ```

Once the database is up, you can run some sample diff operations using the files in `test_data/`:

   ```sh
   $ bin/gh-mysql-conf-diff/gh-mysql-conf-diff test_data/my_sample1.cnf localhost:3306

   Difference found for: CONNECT_TIMEOUT
     my.cnf:    60
     mysqld:    10
   ```

To connect to mysql using the console:

   ```sh
   $ docker-compose exec mysql-database mysql -uroot -p
   ```

and use the password for the test database as defined in [docker-compose.yml](docker-compose.yml).

## Coding Conventions

Our codebase adheres to certain coding conventions. Before you contribute, please make sure to follow them:

1. **Formatting**: We use `gofmt` to automatically format our code. Please make sure to run `gofmt -s -w .` on your code before committing.

2. **Naming**: We prefer short, concise names for local variables and more descriptive names for exported functions and variables. Acronyms should be all uppercase.

3. **Error Handling**: Always check errors and handle them immediately. Do not ignore errors or use panic for normal error handling.

4. **Comments**: Write a comment for every exported function, variable, and type. The comment should start with the name of the thing it's describing.

5. **Packages**: Each package should have a single purpose and provide a clean, simple API. The package name should be a noun, and the functions in the package should be actions on that noun.

6. **Testing**: Write tests for your code. Test functions should be named `TestXxx`, where `Xxx` is the name of the function being tested.

We use `golangci-lint` to maintain the quality of our codebase. We use it in our CI/CD pipeline to catch any issues before they get merged. The linter configuration is at `.golangci.yml`.

## Submitting a pull request

1. [Fork][fork] and clone the repository
1. Make sure the tests pass on your machine: `go test -v ./...`
1. Make sure linter passes on your machine: `golangci-lint run`
1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and make sure the tests and linter still pass
1. Push to your fork and [submit a pull request][pr]
1. Wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Follow the [style guide][style].
- Write tests.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
