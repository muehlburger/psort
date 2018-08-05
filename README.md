# psort

psort sorts photos based on date

[Features](#features) | [Installation](#installation) | [Usage](#usage) | [Examples](#examples) | [Command-line options](#options) | [Configuration](#configuration)

![screenshot](https://raw.githubusercontent.com/muehlburger/psort/master/docs/screenshot.png)

## Features

- **Copy** files from source directories to one target directory

- **Rename** copied files to a date format.

## Installation

```bash
go get github.com/muehlburger/psort
```

## Usage

### Examples

Copy photos:

```bash
psort -s /home/herbert/pic1,/home/herbert/pic1,/home/herbert/pic1 target
```

## Authors

[Herbert Mühlburger](https://github.com/muehlburger) and [contributors](https://github.com/muehlburger/psort/graphs/contributors).

## License

[Apache 2.0 License](LICENSE)

[report-card-image]: https://goreportcard.com/badge/github.com/muehlburger/psort
[report-card-url]: https://goreportcard.com/report/github.com/muehlburger/psort