#ifndef _XSLT_H
#define _XSLT_H

#include <libxslt/xsltutils.h>

struct result {
  int ok;
  int doc_txt_len;
  xmlChar *doc_txt_ptr;
};

int transform(xmlDocPtr style_doc, xmlDocPtr xml_doc, xmlChar **doc_txt_ptr,
              int *doc_txt_len);

struct result xslt(const char *xsl, const char *xml);

#endif
