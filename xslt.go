package xslt

// #cgo LDFLAGS: -L/usr/lib -L/usr/local/lib -lxml2 -lxslt -lz -llzma -lm
// #cgo CFLAGS: -I/usr/include -I/usr/local/include -I/usr/include/libxml2 -I/usr/local/include/libxml2
// #include "xslt.h"
import "C"

import (
	"errors"
	"unsafe"
)

// ErrXSLTFailure is an XSL transformation error
var ErrXSLTFailure = errors.New("XSL transformation failed")

// Transform executes an XSL transformation via Libxslt
func Transform(xsl, xml []byte) ([]byte, error) {

	var (
		cres       C.struct_result
		cxml, cxsl *C.char
		ptr        unsafe.Pointer
	)

	cxsl = C.CString(string(xsl))
	cxml = C.CString(string(xml))
	defer C.free(unsafe.Pointer(cxml))
	defer C.free(unsafe.Pointer(cxsl))

	cres = C.xslt(cxsl, cxml)
	if cres.ok != 0 {
		return nil, ErrXSLTFailure
	}

	ptr = unsafe.Pointer(cres.doc_txt_ptr)
	defer C.free(ptr)

	return C.GoBytes(ptr, cres.doc_txt_len), nil
}
