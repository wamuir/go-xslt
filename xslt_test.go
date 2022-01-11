package xslt_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/wamuir/go-xslt"
)

func TestNewStylesheet(t *testing.T) {
	tests := []struct {
		name    string
		xslFile string
	}{
		{"style1", "testdata/style1.xsl"},
		{"style2", "testdata/style2.xsl"},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			xsl, _ := ioutil.ReadFile(c.xslFile)

			xs, err := xslt.NewStylesheet(xsl)
			if err != nil {
				t.Errorf("got: %v, want: %v", err, nil)
			}

			if xs == nil {
				t.Errorf("got: %v, want: %v", xs, "non-nil")
			}
		})
	}

	errorTests := []struct {
		name   string
		xslStr string
		err    error
	}{
		{"emptyXSL", "", xslt.ErrXSLParseFailure},
	}

	for _, c := range errorTests {
		t.Run(c.name, func(t *testing.T) {
			xs, err := xslt.NewStylesheet([]byte(c.xslStr))
			if xs != nil {
				t.Errorf("got: %v, expected %v", xs, nil)
			}

			if err != c.err {
				t.Errorf("got: %v, expected %v", err, c.err)
			}
		})
	}
}

func TestStylesheetClose(t *testing.T) {
	tests := []struct {
		name    string
		xslFile string
	}{
		{"style1", "testdata/style1.xsl"},
		{"style2", "testdata/style2.xsl"},
		{"nil", ""},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			xs := new(xslt.Stylesheet)
			if len(c.xslFile) > 0 {
				xsl, _ := ioutil.ReadFile(c.xslFile)
				xs, _ = xslt.NewStylesheet(xsl)
			}
			func(xs *xslt.Stylesheet) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("unexpected panic: %v", r)
					}
				}()
				xs.Close()
			}(xs)
		})
	}
}

func TestStylesheetTransform(t *testing.T) {
	tests := []struct {
		name    string
		xmlFile string
		xslFile string
		resFile string
	}{
		{"style1", "testdata/document.xml", "testdata/style1.xsl", "testdata/result1.xml"},
		{"style2", "testdata/document.xml", "testdata/style2.xsl", "testdata/result2.xhtml"},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			xml, _ := ioutil.ReadFile(c.xmlFile)
			xsl, _ := ioutil.ReadFile(c.xslFile)

			xs, _ := xslt.NewStylesheet(xsl)

			got, err := xs.Transform(xml)
			if err != nil {
				t.Errorf("got: %v, want: %v", err, nil)
			}

			want, _ := ioutil.ReadFile(c.resFile)
			if !bytes.Equal(got, want) {
				t.Errorf("got: %s, want: %s", got, want)
			}
		})
	}

	errorTests := []struct {
		name    string
		xmlStr  string
		xslFile string
		err     error
	}{
		{"emptyXML", "", "testdata/style2.xsl", xslt.ErrXSLTFailure},
	}

	for _, c := range errorTests {
		t.Run(c.name, func(t *testing.T) {
			xsl, _ := ioutil.ReadFile(c.xslFile)

			xs, _ := xslt.NewStylesheet(xsl)

			got, err := xs.Transform([]byte(c.xmlStr))
			if got != nil {
				t.Errorf("got: %s, want: %v", got, nil)
			}
			if err != c.err {
				t.Errorf("got: %v, want: %v", err, c.err)
			}
		})
	}
}

func TestStylesheetTransformExslt(t *testing.T) {
	tests := []struct {
		name string
		xml  []byte
		xsl  []byte
		res  []byte
	}{
		{
			"math/min",
			[]byte(`<?xml version="1.0" encoding="UTF-8"?>

<values>
   <value>7</value>
   <value>11</value>
   <value>8</value>
   <value>4</value>
</values>
`),
			[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet
   version="1.0"
   xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
   xmlns:math="http://exslt.org/math"
   extension-element-prefixes="math">
   <xsl:output method="xml" indent="yes" encoding="UTF-8"/>
   <xsl:template match="values">
     <result>
       <xsl:text>Minimum: </xsl:text>
       <xsl:value-of select="math:min(value)" />
     </result>
  </xsl:template>
</xsl:stylesheet>
`),
			[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<result>Minimum: 4</result>
`),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			xs, err := xslt.NewStylesheet(c.xsl)
			if err != nil {
				t.Fatal(err)
			}

			got, err := xs.Transform(c.xml)
			if err != nil {
				t.Fatal(err)
			}

			want := c.res
			if !bytes.Equal(got, want) {
				t.Errorf("got: %s, want: %s", got, want)
			}
		})
	}

}

func BenchmarkStylesheetTransform(b *testing.B) {
	xml, _ := ioutil.ReadFile("testdata/document.xml")
	xsl, _ := ioutil.ReadFile("testdata/style1.xsl")
	xs, _ := xslt.NewStylesheet(xsl)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := xs.Transform(xml); err != nil {
			b.Errorf("got %v, want %v", err, nil)
		}
	}
}
