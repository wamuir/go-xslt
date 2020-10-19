package xslt

// #cgo LDFLAGS: -L/usr/lib -L/usr/local/lib -lxml2 -lxslt -lz -llzma -lm
// #cgo CFLAGS: -I/usr/include -I/usr/local/include -I/usr/include/libxml2 -I/usr/local/include/libxml2
// #include <string.h>
// #include "xslt.h"
import "C"

import (
	"errors"
	"unsafe"
)

// Transformation errors
var (
	ErrXSLTFailure     = errors.New("XSL transformation failed")
	ErrXSLParseFailure = errors.New("Failed to parse XSL")
)

// Stylesheet represents an XSL
type Stylesheet struct {
	ptr C.xsltStylesheetPtr
}

// Close frees memory associated with a stylesheet.  Additional calls
// to Close will be ignored.
func (xs *Stylesheet) Close() {
	if xs.ptr != nil {
		C.free_style(&xs.ptr)
		xs.ptr = nil
	}
}

// Transform applies the receiver to the XML and returns the result of
// an XSL transformation and any errors.  The resulting document may
// be nil (a zero-length and zero-capacity byte slice) in the case of
// an error.
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

// NewStylesheet creates and returns a new stylesheet along with any
// errors.  The resulting stylesheet ay be nil if an error is
// encountered during parsing.  This implementation relies on Libxslt,
// which supports XSLT 1.0.
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
