# BookPrint

[![Docker](../../actions/workflows/docker.yml/badge.svg)](../../actions/workflows/docker.yml)

A Go-based, open-source CLI tool without dependencies that converts HTML files into books.

## ⚙️ Get Started

You'll need [Go](https://go.dev) installed.

### Install

First of all, you need to install `bookprint` locally:

```shell
$ go install stefanco.de/bookprint/cmd/bookprint
```

This will install `bookprint` into `$GOROOT/bin`.

### Run locally

Then you're able to run `bookprint` locally:

```shell
$ bookprint --help
```

## 👨‍💻 Usage

This message is also available when running `$ bookprint --help`.

```text
Usage:
        bookprint [options...] <file>

Options:
        -t, --template-dir <dir>    Path to the directory containing custom templates used for generating the book.
        -o, --output-dir <dir>      Path to the directory where the generated book pages will be stored.
        -s, --static-dir <dir>      Path to the directory with additional files for the book. Copied to output directory.
        -v, --version               Print the version number.
        -h, --help                  Print the help message.

Examples:
        Reading from HTML file:
        $ bookprint --template-dir templates --output-dir out examples/index.html
        > Created book in 'out' directory

        Reading from STDIN:
        $ echo "<html>...</html>" | bookprint --template-dir templates --output-dir out --
        > Created book in 'out' directory
```

## 🔨 Technology

The following technologies, tools and platforms were used during development.

- **Code**: [Go](https://go.dev)
- **CI/CD**: [GitHub Actions](https://github.com/actions)

## 👷‍ Error Found?

Thank you for your message! Please fill out a [bug report](../../issues/new?assignees=&labels=&template=bug_report.md&title=).

## License

This project is licensed under the [European Union Public License 1.2](https://choosealicense.com/licenses/eupl-1.2/).