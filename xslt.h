#ifndef GOXSLT_H
#define GOXSLT_H

#include <libxslt/xsltutils.h>

int apply_style(xsltStylesheetPtr style, const char *xml, const char **params,
                char **xml_txt, size_t *xml_txt_len);

void free_style(xsltStylesheetPtr *style);

int make_style(const char *xsl, xsltStylesheetPtr *style);

int xslt(const char *xsl, const char *xml, const char **params, char **xml_txt,
         size_t *xml_txt_len);

void init_exslt();

const char **make_param_array(int num_tuples);

void set_param(char **a, char *n, char *v, int t);

void free_param_array(char **a, int num_tuples);

#endif
