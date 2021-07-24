package xslt

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/document.xml
var document []byte

//go:embed testdata/style1.xsl
var style1 []byte

//go:embed testdata/style2.xsl
var style2 []byte

//go:embed testdata/result1.xml
var result1 []byte

//go:embed testdata/result2.xhtml
var result2 []byte

func TestNewStylesheet(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs1 == nil {
		t.Errorf("xs1 = %v, want non-nil", xs1)
	}
	defer xs1.Close()

	xs2, err := NewStylesheet(style2)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs2 == nil {
		t.Errorf("xs2 = %v, want non-nil", xs2)
	}
	defer xs2.Close()

	xs3, err := NewStylesheet(nil)
	if err != ErrXSLParseFailure {
		t.Errorf("err = %v; want %v", err, ErrXSLParseFailure)
	}
	if xs3 != nil {
		t.Errorf("xs3 = %v, want %v", xs3, nil)
	}
}

func TestStylesheetClose(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs1 == nil {
		t.Errorf("xs1 = %v, want non-nil", xs1)
	}
	xs1.Close()

	xs2 := Stylesheet{}
	xs2.Close()
}

func TestStylesheetTransform(t *testing.T) {

	xs1, err := NewStylesheet(style1)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs1 == nil {
		t.Errorf("xs1 = %v, want non-nil", xs1)
	}
	defer xs1.Close()

	res1, err := xs1.Transform(document)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if !bytes.Equal(result1, res1) {
		t.FailNow()
	}

	xs2, err := NewStylesheet(style2)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs2 == nil {
		t.Errorf("xs2 = %v, want non-nil", xs2)
	}
	defer xs2.Close()

	res2, err := xs2.Transform(document)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if !bytes.Equal(result2, res2) {
		t.FailNow()
	}

	xs3, err := NewStylesheet(style1)
	if err != nil {
		t.Errorf("err = %v; want %v", err, nil)
	}
	if xs3 == nil {
		t.Errorf("xs3 = %v, want non-nil", xs3)
	}
	defer xs3.Close()

	res3, err := xs3.Transform(nil)
	if err != ErrXSLTFailure {
		t.Errorf("err = %v; want %v", err, ErrXSLTFailure)
	}
	if res3 != nil {
		t.Errorf("xs3 = %v, want %v", res3, nil)
	}
}
