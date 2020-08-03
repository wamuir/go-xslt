#ifndef GOXSLT_H
#define GOXSLT_H

#include <libxslt/xsltutils.h>

int apply_style(xsltStylesheetPtr style, const char *xml, char **xml_txt,
                size_t *xml_txt_len);

void free_style(xsltStylesheetPtr *style);

int make_style(const char *xsl, xsltStylesheetPtr *style);

int xslt(const char *xsl, const char *xml, char **xml_txt, size_t *xml_txt_len);

#endif
