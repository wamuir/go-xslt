package xslt

/*
#cgo LDFLAGS: -lxml2 -lxslt -lexslt -lz -llzma -lm
#cgo CFLAGS: -I/usr/include -I/usr/include/libxml2
#cgo freebsd LDFLAGS: -L/usr/local/lib
#cgo freebsd CFLAGS: -I/usr/local/include -I/usr/local/include/libxml2
#include <string.h>
#include "xslt.h"
*/
import "C"

import (
	"errors"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// Package errors.
var (
	ErrMixedQuotes     = errors.New("unable to quote parameter value")
	ErrUTF8Validation  = errors.New("input failed utf-8 validation")
	ErrXSLTFailure     = errors.New("xsl transformation failed")
	ErrXSLParseFailure = errors.New("failed to parse xsl")
)

func init() {
	C.init_exslt()
}

// Parameter is a parameter to be passed to a stylesheet.
type Parameter interface {
	name() (*C.char, error)
	value() (*C.char, error)
}

// StringParameter is a stylesheet parameter consisting of a name/value pair,
// where name and value are UTF-8 strings.
type StringParameter struct {
	Name  string
	Value string
}

var _ Parameter = StringParameter{}

func (p StringParameter) name() (*C.char, error) {
	if ok := utf8.ValidString(p.Name); !ok {
		return nil, ErrUTF8Validation
	}

	return C.CString(p.Name), nil
}

func (p StringParameter) value() (*C.char, error) {
	var (
		r rune
		s string
	)

	const (
		doublequote = rune(0x0022)
		singlequote = rune(0x0027)
	)

	if strings.ContainsRune(p.Value, doublequote) {
		if strings.ContainsRune(p.Value, singlequote) {
			return nil, ErrMixedQuotes
		}
		r = singlequote
	} else {
		r = doublequote
	}

	s = string(r) + p.Value + string(r)
	if ok := utf8.ValidString(s); !ok {
		return nil, ErrUTF8Validation
	}

	return C.CString(s), nil
}

// XPathParameter is a stylesheet parameter consisting of a name/value pair,
// where name is a QName or a UTF-8 string of the form {URI}NCName and value is
// a UTF-8 XPath expression. A quoted value (single or double) will be treated
// as a string rather than as an XPath expression, however the use of
// StringParameter is preferable when passing string parameters to a
// stylesheet.
type XPathParameter struct {
	Name  string
	Value string
}

var _ Parameter = XPathParameter{}

func (p XPathParameter) name() (*C.char, error) {
	if ok := utf8.ValidString(p.Name); !ok {
		return nil, ErrUTF8Validation
	}

	return C.CString(p.Name), nil
}

func (p XPathParameter) value() (*C.char, error) {
	if ok := utf8.ValidString(p.Value); !ok {
		return nil, ErrUTF8Validation
	}

	return C.CString(p.Value), nil
}

// Stylesheet represents an xsl stylesheet.
type Stylesheet struct {
	ptr C.xsltStylesheetPtr
}

// Close frees memory associated with a stylesheet.  Additional calls to Close
// will be ignored.
func (xs *Stylesheet) Close() {
	if xs.ptr != nil {
		C.free_style(&xs.ptr)
		xs.ptr = nil
	}
}

// Transform applies receiver stylesheet xs to xml and returns the result of an
// xsl transformation and any error.  The resulting document may be nil (a
// zero-length and zero-capacity byte slice) in the case of an error.
func (xs *Stylesheet) Transform(xml []byte, params ...Parameter) ([]byte, error) {
	var (
		cxml *C.char
		cout *C.char
		cpar **C.char
		ret  C.int
		size C.size_t
	)

	cxml = C.CString(string(xml))
	defer C.free(unsafe.Pointer(cxml))

	cpar = C.make_param_array(C.int(len(params)))
	for i, p := range params {
		cname, err := p.name()
		if err != nil {
			return nil, err
		}

		cval, err := p.value()
		if err != nil {
			return nil, err
		}

		C.set_param(cpar, cname, cval, C.int(i))
	}
	defer C.free_param_array(cpar, C.int(len(params)))

	ret = C.apply_style(xs.ptr, cxml, cpar, &cout, &size)

	ptr := unsafe.Pointer(cout)
	defer C.free(ptr)

	if ret != 0 {
		return nil, ErrXSLTFailure
	}

	return C.GoBytes(ptr, C.int(size)), nil
}

// NewStylesheet creates and returns new stylesheet xs along with any error.
// The resulting stylesheet may be nil if an error is encountered during
// parsing.  This implementation relies on Libxslt, which supports XSLT 1.0.
func NewStylesheet(xsl []byte) (*Stylesheet, error) {
	var (
		cxsl *C.char
		cssp C.xsltStylesheetPtr
		ret  C.int
	)

	cxsl = C.CString(string(xsl))
	defer C.free(unsafe.Pointer(cxsl))

	ret = C.make_style(cxsl, &cssp)
	if ret != 0 {
		return nil, ErrXSLParseFailure
	}

	return &Stylesheet{ptr: cssp}, nil
}
