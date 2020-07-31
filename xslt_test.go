package xslt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var xml = []byte(
	`<?xml version="1.0" encoding="utf-8" ?>
<rodents>
<rodent>gopher</rodent>
</rodents>`,
)

var xsl = []byte(
	`<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:template match="/">
<xsl:for-each select="rodents/rodent">
<h1><xsl:value-of select="."/></h1>
</xsl:for-each>
</xsl:template>
</xsl:stylesheet>`,
)

var exp = []byte(
	`<?xml version="1.0"?>
<h1>gopher</h1>
`,
)

func TestTransform_A(t *testing.T) {

	out, err := Transform(xsl, xml)
	assert.Nil(t, err)

	assert.Equal(t, exp, out)
}

// Empty XML
func TestTransform_B(t *testing.T) {

	_, err := Transform(xsl, make([]byte, 0))
	assert.EqualError(t, ErrXSLTFailure, err.Error())
}

// Empty XSL
func TestTransform_C(t *testing.T) {

	_, err := Transform(make([]byte, 0), xml)
	assert.EqualError(t, ErrXSLTFailure, err.Error())
}

// Invalid XML
func TestTransform_D(t *testing.T) {

	_, err := Transform(xsl, xml[0:len(xml)-1])
	assert.EqualError(t, ErrXSLTFailure, err.Error())
}

// Invalid XSL
func TestTransform_E(t *testing.T) {

	_, err := Transform(xsl[0:len(xsl)-1], xml)
	assert.EqualError(t, ErrXSLTFailure, err.Error())
}
