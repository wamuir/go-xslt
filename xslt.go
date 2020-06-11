package xslt

/*
#cgo LDFLAGS: -L/usr/lib -L/usr/local/lib -lxml2 -lxslt -lz -llzma -lm
#cgo CFLAGS: -I/usr/include -I/usr/local/include -I/usr/include/libxml2 -I/usr/local/include/libxml2
#include <string.h>
#include <libxml/xmlreader.h>
#include <libxslt/transform.h>
#include <libxslt/xsltutils.h>

struct Result {
	int         ok;
	const char *r;
};

struct Result transform(const char * xsl, const char * xml) {

	struct Result res = {0, NULL};

	xmlDocPtr style = xmlParseMemory(xsl, strlen(xsl));
	xmlDocPtr doc = xmlParseMemory(xml, strlen(xml));

	xmlErrorPtr e = xmlGetLastError();

	if (e == NULL) {

		xsltStylesheetPtr transformer = xsltParseStylesheetDoc(style);
		if ((transformer != NULL) && (transformer->errors == 0)) {
			xmlDocPtr newDoc = xsltApplyStylesheet(
				transformer, doc, NULL
			);
			if (newDoc != NULL) {
				int buffersize;
				xmlChar *result;
       				xsltSaveResultToString(
					&result, &buffersize, newDoc, transformer
				);
				res.r = (const char *) result;
				res.ok = 1;
				xmlFreeDoc(newDoc);
			}
			xsltFreeStylesheet(transformer);
		}
		xsltCleanupGlobals();

	} else {

		xmlFreeDoc(style);

	}

	xmlFreeDoc(doc);
	xmlCleanupParser();

	return res;
}

*/
import "C"

import (
	"errors"
	"unsafe"
)

// Execute an XSL transformation via Libxslt
func Transform(xsl, xml []byte) ([]byte, error) {

	var result []byte

	var cxsl *C.char = C.CString(string(xsl))
	var cxml *C.char = C.CString(string(xml))

	res := C.transform(cxsl, cxml)
	C.free(unsafe.Pointer(cxml))
	C.free(unsafe.Pointer(cxsl))

	if res.ok == 1 {
		result = []byte(C.GoString(res.r))
		defer C.free(unsafe.Pointer(res.r))
	} else {
		err := errors.New("error encountered in XSL transform")
		return make([]byte, 0), err
	}

	return result, nil
}
