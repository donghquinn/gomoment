# GoMoment

**English** | [한국어](./README.ko.md)

[![Go Version](https://img.shields.io/badge/go-1.12+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/donghquinn/gomoment.svg)](https://pkg.go.dev/github.com/donghquinn/gomoment)

A Go library for parsing, formatting, and converting dates between timezones with [moment.js](https://momentjs.com/)-style format tokens, so you can write `YYYY-MM-DD` instead of Go's reference time `2006-01-02`.

```go
moment, _ := gomoment.NewMoment("2023-05-15T14:30:45Z")
kst, _ := moment.TZ("KST")
fmt.Println(kst.Must("YYYY-MM-DD HH:mm Z")) // 2023-05-15 23:30 +09:00
```

## Features

- **moment.js-style tokens** — `YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`, `SSS`, `A`, `Z`, ...
- **Timezone conversion** — `TZ("KST")`, `TZ("America/New_York")`, `UTC()`, `Local()`
- **Flexible parsing** — common date string formats are detected automatically
- **Strict validation** — unknown tokens return an error instead of silently producing garbage
- **Zero dependencies** — standard library only

## Installation

```bash
go get -u github.com/donghquinn/gomoment
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/donghquinn/gomoment"
)

func main() {
    now, err := gomoment.NewMoment()
    if err != nil {
        panic(err)
    }

    formatted, err := now.Format("YYYY-MM-DD HH:mm:ss")
    if err != nil {
        panic(err)
    }

    fmt.Println(formatted) // e.g. 2025-08-06 14:30:45
}
```

## Usage

### Creating Moments

#### Current Time
```go
moment, err := gomoment.NewMoment()
// or
moment := gomoment.Now()
```

#### From String
```go
moment, err := gomoment.NewMoment("2025-08-06 14:30:45")
// or
moment, err := gomoment.Parse("2025-08-06 14:30:45")
```

#### From time.Time
```go
t := time.Now()
moment, err := gomoment.NewMoment(t)
```

Any other argument type returns an error.

### Formatting

| Token | Description | Example |
|-------|-------------|---------|
| `YYYY` | 4-digit year | 2025 |
| `YY` | 2-digit year | 25 |
| `MM` | Month (01-12) | 08 |
| `M` | Month (1-12) | 8 |
| `DD` | Day of month (01-31) | 06 |
| `D` | Day of month (1-31) | 6 |
| `HH` | Hour, 24-hour clock (00-23) | 09 |
| `H` | Hour, 24-hour clock (0-23) | 9 |
| `hh` | Hour, 12-hour clock (01-12) | 02 |
| `h` | Hour, 12-hour clock (1-12) | 2 |
| `mm` | Minutes (00-59) | 05 |
| `m` | Minutes (0-59) | 5 |
| `ss` | Seconds (00-59) | 07 |
| `s` | Seconds (0-59) | 7 |
| `SSS` | Milliseconds (000-999) | 123 |
| `A` | AM/PM | PM |
| `a` | am/pm | pm |
| `ZZ` | UTC offset | +0900 |
| `Z` | UTC offset with colon | +09:00 |

Formatting rules:

- A token is a run of one repeated letter. A run that is not in the table (for example `MMM` or `DDDD`) is an error.
- Text inside square brackets is printed as-is: `[Today is] YYYY-MM-DD` → `Today is 2025-08-06`.
- Characters other than ASCII letters (digits, punctuation, spaces, Korean text, ...) are printed as-is.

#### Examples

```go
moment, _ := gomoment.NewMoment("2025-08-06 14:30:45")

fmt.Println(moment.Must("YYYY-MM-DD"))          // 2025-08-06
fmt.Println(moment.Must("MM/DD/YYYY"))          // 08/06/2025
fmt.Println(moment.Must("YYYY-MM-DD HH:mm:ss")) // 2025-08-06 14:30:45
fmt.Println(moment.Must("h:mm A"))              // 2:30 PM
fmt.Println(moment.Must("[Date:] YYYY.MM.DD"))  // Date: 2025.08.06
fmt.Println(moment.Must("YYYY년 M월 D일"))        // 2025년 8월 6일
```

### Must

`Must` works like `Format` but panics on an invalid format string. Use it only with format strings you know are valid, such as constants:

```go
moment := gomoment.Now()
formatted := moment.Must("YYYY-MM-DD HH:mm:ss")
```

### Supported Input Formats

`NewMoment(string)` and `Parse` try these layouts in order (shown in Go reference-time notation):

- `2006-01-02`
- `2006/01/02`
- `2006-01-02 15:04:05`
- `2006/01/02 15:04:05`
- `2006-01-02T15:04:05`
- `2006-01-02T15:04:05Z`
- `2006-01-02T15:04:05-07:00`
- `01/02/2006`
- `01-02-2006`
- `01/02/2006 15:04:05`
- `01-02-2006 15:04:05`
- `15:04:05`

Strings without an explicit offset are interpreted as UTC. Use `TZ` to convert the result to another timezone.

### Getting the underlying time.Time

```go
moment, _ := gomoment.NewMoment()
t := moment.Time() // time.Time
```

### Timezone Support

Every conversion returns a new `Moment`; the original is never modified.

#### Convert to a Specific Timezone

```go
moment, _ := gomoment.NewMoment("2023-05-15T14:30:45Z")

kstMoment, _ := moment.TZ("KST")
fmt.Println(kstMoment.Must("YYYY-MM-DD HH:mm")) // 2023-05-15 23:30

estMoment, _ := moment.TZ("EST")
fmt.Println(estMoment.Must("YYYY-MM-DD HH:mm")) // 2023-05-15 10:30 (EDT in May)

// IANA timezone names also work
nyMoment, _ := moment.TZ("America/New_York")
tokyoMoment, _ := moment.TZ("Asia/Tokyo")
```

#### Supported Timezone Abbreviations

Abbreviations are case-insensitive. Each one maps to an IANA location, so daylight saving time is applied where that location observes it (for example, `EST` follows New York and becomes EDT in summer).

| Abbreviation | IANA Name | Description |
|--------------|-----------|-------------|
| `UTC` | UTC | Coordinated Universal Time |
| `GMT` | GMT | Greenwich Mean Time |
| `KST` | Asia/Seoul | Korea Standard Time |
| `JST` | Asia/Tokyo | Japan Standard Time |
| `EST` | America/New_York | Eastern Time |
| `PST` | America/Los_Angeles | Pacific Time |
| `CST` | America/Chicago | Central Time |
| `MST` | America/Denver | Mountain Time |
| `CET` | Europe/Paris | Central European Time |
| `IST` | Asia/Kolkata | India Standard Time |

#### Other Timezone Helpers

```go
moment, _ := gomoment.NewMoment()

utcMoment := moment.UTC()     // convert to UTC
localMoment := moment.Local() // convert to the local timezone

name, offset := moment.Zone() // e.g. "KST", 32400
offset = moment.Offset()      // seconds east of UTC
```

## Error Handling

```go
// Invalid date string
_, err := gomoment.NewMoment("invalid-date")
fmt.Println(err)
// unable to parse date string "invalid-date" (supported formats: YYYY-MM-DD, YYYY/MM/DD, RFC3339, etc.)

// Unsupported format tokens
moment := gomoment.Now()
_, err = moment.Format("YYYY-QQ-MMM")
fmt.Println(err)
// unsupported format tokens: [QQ MMM] (supported tokens: YYYY, YY, MM, M, DD, D, HH, H, hh, h, mm, m, ss, s, SSS, A, a, ZZ, Z)

// Unknown timezone
_, err = moment.TZ("Mars/Olympus")
fmt.Println(err) // invalid timezone: Mars/Olympus
```

## Performance

`Format` writes the output in a single pass over the format string. It does not build an intermediate Go layout, so it typically makes 1–2 allocations per call.

```
goos: darwin
goarch: arm64
cpu: Apple M5
```

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `Format("YYYY-MM-DD")` | 54.6 | 16 | 1 |
| `Format("YYYY-MM-DD HH:mm:ss Z")` | 130.7 | 80 | 2 |
| stdlib `time.Format("2006-01-02 15:04:05 -07:00")` | 66.0 | 32 | 1 |
| `NewMoment()` | 34.6 | 24 | 1 |
| `NewMoment("2023-05-15 14:30:45")` | 169.0 | 312 | 9 |
| `Now()` | 27.5 | 0 | 0 |

`Format` costs roughly 2× the standard library's `time.Format` for an equivalent layout. In exchange, you write moment.js tokens instead of Go's reference time.

Run the benchmarks yourself:

```bash
go test -run '^$' -bench . -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

```bash
go vet ./...
go test -race ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Inspiration

This library is inspired by the excellent [moment.js](https://momentjs.com/) library for JavaScript.
