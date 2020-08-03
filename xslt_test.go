package xslt

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	document = mustReadFile("testdata/document.xml")
	style1   = mustReadFile("testdata/style1.xsl")
	style2   = mustReadFile("testdata/style2.xsl")
	result1  = mustReadFile("testdata/result1.xml")
	result2  = mustReadFile("testdata/result2.xhtml")
)

func mustReadFile(f string) []byte {

	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return b
}

func TestNewStylesheet(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	assert.NotNil(t, xs1)
	assert.NoError(t, err)
	defer xs1.Close()

	xs2, err := NewStylesheet(style2)
	assert.NotNil(t, xs2)
	assert.NoError(t, err)
	defer xs2.Close()

	xs3, err := NewStylesheet(nil)
	assert.Nil(t, xs3)
	assert.Equal(t, ErrXSLParseFailure, err)
}

func TestStylesheetClose(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	assert.NotNil(t, xs1)
	assert.NoError(t, err)
	assert.NotPanics(t, xs1.Close)

	xs2 := Stylesheet{}
	assert.NotPanics(t, xs2.Close)
}

func TestStylesheetTransform(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	assert.NotNil(t, xs1)
	assert.NoError(t, err)
	defer xs1.Close()

	res1, err := xs1.Transform(document)
	assert.Equal(t, result1, res1)
	assert.NoError(t, err)

	xs2, err := NewStylesheet(style2)
	assert.NotNil(t, xs2)
	assert.NoError(t, err)
	defer xs2.Close()

	res2, err := xs2.Transform(document)
	assert.Equal(t, result2, res2)
	assert.NoError(t, err)

	xs3, err := NewStylesheet(style1)
	assert.NotNil(t, xs3)
	assert.NoError(t, err)
	defer xs3.Close()

	res3, err := xs3.Transform(nil)
	assert.Nil(t, res3)
	assert.Equal(t, ErrXSLTFailure, err)
}
