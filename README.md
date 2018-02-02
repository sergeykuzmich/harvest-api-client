Go Harvest API
==============

[![GoDoc](https://godoc.org/github.com/adlio/harvest?status.svg)](http://godoc.org/github.com/adlio/harvest)
[![Build Status](https://travis-ci.org/adlio/harvest.svg)](https://travis-ci.org/adlio/harvest)
[![Coverage Status](https://coveralls.io/repos/github/adlio/harvest/badge.svg?branch=master)](https://coveralls.io/github/adlio/harvest?branch=master)

A #golang package to access the [Harvest API](https://help.getharvest.com/api-v2/).


## Installation

The Go Harvest API has been tested compatible with Go 1.8 on up. Its only dependency is
the `github.com/pkg/errors` package. It otherwise relies only on the Go standard library.

```
go get github.com/adlio/harvest
```

## Basic Usage

All interaction starts with a `harvest.API`. Create one with your account ID and token:

```Go
client := harvest.NewTokenAPI("ACCOUNTID", "TOKEN")
```
