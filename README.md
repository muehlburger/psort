# sortphotos

sortphotos sorts photos based on EXIF information and flattens the directory structure.

[Features](#features) | [Installation](#installation) | [Usage](#usage) | [Examples](#examples) | [Command-line options](#options) | [Configuration](#configuration)

![screenshot](https://raw.githubusercontent.com/muehlburger/sortphotos/master/docs/screenshot.png)

## Features

- **Copy** files from source directories to one target directory

- **Rename** copied files to a date format.

## Installation

```bash
go get -u github.com/muehlburger/sortphotos
```

## Usage

### Examples

Copy photos:

```bash
sortphotos photos-to-sort photos-1-sorted
```

## Authors

[Herbert Mühlburger](https://github.com/muehlburger) and [contributors](https://github.com/muehlburger/sortphotos/graphs/contributors).

## License

[Apache 2.0 License](LICENSE)

[report-card-image]: https://goreportcard.com/badge/github.com/muehlburger/sortphotos
[report-card-url]: https://goreportcard.com/report/github.com/muehlburger/sortphotos