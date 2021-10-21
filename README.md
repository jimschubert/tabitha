# tabitha

tabitha is a no-frills tabular formatter for the terminal.

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/tabitha)
[![Go Build](https://github.com/jimschubert/tabitha/actions/workflows/build.yml/badge.svg)](https://github.com/jimschubert/tabitha/actions/workflows/build.yml)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/jimschubert/tabitha?color=orange&label=Docker%20Image%20Size)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/tabitha)](https://goreportcard.com/report/github.com/jimschubert/tabitha)

## Features

* Supports padding output to the longest _display_ text, honoring ANSI colors
* Simple API to add a single header and 0..n lines
* Add "spacer" lines anywhere in the table. Character output is automatically calculated based on width of other text.
* Customize line separator, padding character, disable padding or honoring ansi

## Build/Test

```shell
go test -v -race -cover ./...
```

## Example

The following example demonstrates usage (sans error handling).
```
tt := tabitha.NewWriter()
tt.Header("First", "Second", "Third", "Fourth")
tt.SpacerLine()
tt.AddLine("I'm first", "I'm second", "I'm third", "I'm fourth")
tt.SpacerLine()
tt.WriteTo(os.Stdout)
```

## Why not text/tabwriter?

This is a little different from `text/tabwriter`; a main difference being that tabwriter assumes all runes are the same size and includes ANSI codes in width calculations while `tabitha` does not.

In `tabitha`, rune width is evaluated using `utf8.RuneCountInString`. When ANSI support is enabled, `tabitha` will extract non-ANSI text (the "displayable text") using regex. Since `tabitha` collects all tabular text before writing out to a target `io.Writer`, it is not expected to perform as well as `tabwriter`.

## License

This project is [licensed](./LICENSE) under Apache 2.0.
