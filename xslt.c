#include <libxslt/transform.h>
#include <libxslt/xsltutils.h>
#include <string.h>
#include <xslt.h>

/*
 * Function: transform
 * ----------------------------
 *  Transforms an xml document using an xsl stylesheet
 *
 *   style_doc:    parsed xsl stylesheet
 *   xml_doc:      parsed xml to be transformed
 *   doc_txt_ptr:  output from the transform (unsigned char)
 *   doc_txt_len:  length of output from the transform (int)
 *
 *  returns: 0 if the transform is sucessful or -1 in case of error
 */
int transform(xmlDocPtr style_doc, xmlDocPtr xml_doc, xmlChar **doc_txt_ptr,
              int *doc_txt_len) {

  xmlOutputBufferPtr buf;
  xmlDocPtr result;
  xsltStylesheetPtr style;

  if (!(style = xsltParseStylesheetDoc(style_doc)) || (style->errors)) {
    return -1;
  }

  if (!(result = xsltApplyStylesheet(style, xml_doc, NULL))) {
    return -1;
  }

  return xsltSaveResultToString(doc_txt_ptr, doc_txt_len, result, style);
}

/*
 * Function: xslt
 * ----------------------------
 *  Returns the square of the largest
 *
 *   xsl: the stylesheet to be used
 *   xml: the document to transform
 *
 *  returns: result struct { int OK; xmlChar *do_txt_ptr; int doc_txt_len }
 */
struct result xslt(const char *xsl, const char *xml) {

  xmlDocPtr style, doc;
  struct result res;

  style = xmlParseMemory(xsl, strlen(xsl));
  doc = xmlParseMemory(xml, strlen(xml));

  if (xmlGetLastError()) {
    xmlFreeDoc(style);
    xmlCleanupParser();
    res.ok = -1;
    return res;
  }

  res.ok = transform(style, doc, &res.doc_txt_ptr, &res.doc_txt_len);

  xmlFreeDoc(style);
  xmlFreeDoc(doc);
  xmlCleanupParser();
  xsltCleanupGlobals();

  return res;
}
