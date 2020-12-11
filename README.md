go-xslt
=====

[![GoDoc Reference](https://godoc.org/github.com/wamuir/go-xslt?status.svg)](http://godoc.org/github.com/wamuir/go-xslt)
[![Build Status](https://travis-ci.com/wamuir/go-xslt.svg?branch=master)](https://travis-ci.com/wamuir/go-xslt)
[![codecov](https://codecov.io/gh/wamuir/go-xslt/branch/master/graph/badge.svg)](https://codecov.io/gh/wamuir/go-xslt)
[![Go Report Card](https://goreportcard.com/badge/github.com/wamuir/go-xslt)](https://goreportcard.com/report/github.com/wamuir/go-xslt)

# Description

`go-xslt` is a Go module that performs basic XSLT 1.0 transformations via Libxslt.

# Installation

You'll need the development kits for libxml2 and libxslt.  Install these
via your package manager. For instance, if using `apt` then:

    sudo apt install libxml2-dev libxslt-dev

This module can be installed with the `go get` command:

    go get -u github.com/wamuir/go-xslt


# Usage

```go

  // style is an XSLT 1.0 stylesheet, as []byte.
  xs, err := xslt.NewStylesheet(style)
  if err != nil {
      panic(err)
  }
  defer xs.Close()

  // doc is an XML document to be transformed and res is the result of
  // the XSL transformation, both as []byte. 
  res, err := xs.Transform(doc)
  if err != nil {
      panic(err)
  }

```
