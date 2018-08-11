# sortphotos

sortphotos sorts photos based on EXIF information and flattens the directory structure.

[Features](#features) | [Installation](#installation) | [Usage](#usage) | [Examples](#examples) | [Command-line options](#options) | [Configuration](#configuration)

![screenshot](https://raw.githubusercontent.com/muehlburger/sortphotos/master/docs/screenshot.png)

## Features

- Supports multiple source directories
- **Flattens** source directory structure
- **Copies** files from source directories to one target directory
- **Renames** copied files using EXIF information to the following format: YYYY-MM-DD-HH-MM-SS.[jpg|jpeg|png] -> e.g. 2018-04-20_193419.jpg
- **Supports** jpg, jpeg and png

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