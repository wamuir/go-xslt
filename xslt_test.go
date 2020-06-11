package xslt

import "testing"

const (
	xml = `<?xml version="1.0" encoding="utf-8" ?>
<rodents>
<rodent>gopher</rodent>
</rodents>`

	xsl = `<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:template match="/">
<xsl:for-each select="rodents/rodent">
<h1><xsl:value-of select="."/></h1>
</xsl:for-each>
</xsl:template>
</xsl:stylesheet>
`

	exp = `<?xml version="1.0"?>
<h1>gopher</h1>
`
)

func TestTransform(t *testing.T) {

	html, err := Transform([]byte(xsl), []byte(xml))
	if err != nil {
		t.Error(err)
	}

	if string(html) != exp {
		t.Errorf(
			"Transform was incorrect, got: %s, expect: %s.",
			string(html),
			exp,
		)
	}
}
