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
	"unsafe"
)

// Package errors.
var (
	ErrXSLTFailure     = errors.New("xsl transformation failed")
	ErrXSLParseFailure = errors.New("failed to parse xsl")
)

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
func (xs *Stylesheet) Transform(xml []byte) ([]byte, error) {
	var (
		cxml *C.char
		cout *C.char
		ret  C.int
		size C.size_t
	)

	cxml = C.CString(string(xml))
	defer C.free(unsafe.Pointer(cxml))

	ret = C.apply_style(xs.ptr, cxml, &cout, &size)
	if ret != 0 {
		defer C.free(unsafe.Pointer(cout))
		return nil, ErrXSLTFailure
	}

	ptr := unsafe.Pointer(cout)
	defer C.free(ptr)

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

// InitExslt enables exsl namespace. Call this once at program start.
func InitExslt() {
	C.init_exslt()
}
