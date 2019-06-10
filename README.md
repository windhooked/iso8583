# ISO 8583

[![GoDoc](https://godoc.org/github.com/rvflash/iso8583?status.svg)](https://godoc.org/github.com/rvflash/iso8583)
[![Build Status](https://img.shields.io/travis/rvflash/iso8583.svg)](https://travis-ci.org/rvflash/iso8583)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/iso8583.svg)](http://codecov.io/github/rvflash/iso8583?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/iso8583)](https://goreportcard.com/report/github.com/rvflash/iso8583)

The package `iso8583` implements encoding and decoding of message as defined in ISO 8583.
This specification describes the Financial Transaction Message Format.


An ISO 8583 message is structured as follows:

    Message header.
    Message Type Indicator (MTI).
    One or more bitmaps indicating which data elements are present in the message.
    Data elements, or fields.


### Installation
    
To install it, you need to install Go and set your Go workspace first.
Then, download and install it:

```bash
$ go get -u github.com/rvflash/iso8583
```    
Import it in your code:
    
```go
import "github.com/rvflash/iso8583"
```

### Prerequisite

`iso8583` uses the Go modules that required Go 1.11 or later.
