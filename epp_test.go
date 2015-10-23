package eppgo_test

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/rafaeljusto/eppgo"
)

func TestHelloGeneratingXML(t *testing.T) {
	expected := `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello></hello></epp>`

	hello := eppgo.EPP{
		Child: eppgo.Hello{},
	}

	output, err := xml.Marshal(hello)
	if err != nil {
		t.Fatalf("failed to create the XML for Hello, details: %s", err)
	}

	if string(output) != expected {
		t.Errorf("unexpected Hello XML, expected “%s” and got “%s”", expected, string(output))
	}
}

func TestHelloParsingXML(t *testing.T) {
	expected := eppgo.EPP{
		XMLName: xml.Name{
			Space: "urn:ietf:params:xml:ns:epp-1.0",
			Local: "epp",
		},
		Child: eppgo.Hello{
			XMLName: xml.Name{
				Space: "urn:ietf:params:xml:ns:epp-1.0",
				Local: "hello",
			},
		},
	}

	input := []byte(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello></hello></epp>`)

	var epp eppgo.EPP
	if err := xml.Unmarshal(input, &epp); err != nil {
		t.Fatalf("failed to parse XML, details: %s", err)
	}

	if !reflect.DeepEqual(expected, epp) {
		t.Errorf("unexpected XML parsing, expected “%v” and got “%v”", expected, epp)
	}
}

func TestHelloParsingXMLSecondLevel(t *testing.T) {
	expectedError := eppgo.Error{
		Code:      eppgo.ErrorCodeUnknownElement,
		Reference: "foo",
	}

	input := []byte(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><foo><hello></hello></foo></epp>`)

	var epp eppgo.EPP
	err := xml.Unmarshal(input, &epp)

	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("unexpected error, expected “%v” and got “%v”", expectedError, err)
	}
}
