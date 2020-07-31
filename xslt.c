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

  int ok;
  xmlDocPtr result;
  xsltStylesheetPtr style;

  if (!(style = xsltParseStylesheetDoc(style_doc)) || (style->errors)) {
    return -1;
  }

  if (!(result = xsltApplyStylesheet(style, xml_doc, NULL))) {
    return -1;
  }

  ok = xsltSaveResultToString(doc_txt_ptr, doc_txt_len, result, style);

  xmlFreeDoc(result);
  xmlFree(style);

  return ok;
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

  xmlDocPtr style_doc, xml_doc;
  struct result res;

  style_doc = xmlParseMemory(xsl, strlen(xsl));
  xml_doc = xmlParseMemory(xml, strlen(xml));

  res.ok = (xmlGetLastError()) ? -1
                               : transform(style_doc, xml_doc, &res.doc_txt_ptr,
                                           &res.doc_txt_len);

  (xml_doc) ? xmlFreeDoc(xml_doc) : NULL;
  (style_doc) ? xmlFreeDoc(style_doc) : NULL;
  xmlCleanupParser();
  xsltCleanupGlobals();

  return res;
}
