package eppgo

import (
	"encoding/xml"
	"io"
	"reflect"
	"strings"
)

var (
	registeredChildren map[string]interface{}
)

func init() {
	registeredChildren = make(map[string]interface{})
	RegisterChild(Hello{})
}

// RegisterChild register possible children for the EPP element.
func RegisterChild(obj interface{}) error {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return Error{
			Code:      ErrorCodeRegisteringWrongType,
			Reference: value.Kind().String(),
		}
	}

	name := value.Type().Name()

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.Name == "XMLName" {
			xmlData := field.Tag.Get("xml")
			xmlDataParts := strings.Split(xmlData, " ")

			if len(xmlDataParts) == 1 {
				name = xmlDataParts[0]
			} else if len(xmlDataParts) == 2 {
				name = xmlDataParts[1]
			}
		}
	}

	registeredChildren[name] = obj
	return nil
}

// EPP is the protocol identification. This type identifies the start of an EPP
// protocol element and the namespace used within the protocol.
type EPP struct {
	// XMLName as described in https://golang.org/pkg/encoding/xml/#Unmarshal the
	// attribute is necessary to fill the namespace attribute.
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Child   interface{}
}

// UnmarshalXML unmarshal an EPP XML element.
func (e *EPP) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	e.XMLName = start.Name

	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		switch v := t.(type) {
		case xml.StartElement:
			if obj, ok := registeredChildren[v.Name.Local]; ok {
				// copy the registered child
				value := reflect.ValueOf(obj)
				newValue := reflect.New(value.Type())
				newObj := newValue.Interface()

				if err := d.DecodeElement(newObj, &v); err != nil {
					return err
				}

				// don't need to store the pointer of the child in the EPP type
				e.Child = newValue.Elem().Interface()

			} else {
				return Error{
					Code:      ErrorCodeUnknownElement,
					Reference: v.Name.Local,
				}
			}
		}
	}

	return nil
}

// Hello EPP MAY be carried over both connection-oriented and connection-less
// transport protocols.  An EPP client MAY request a <greeting> from an EPP
// server at any time between a successful <login> command and a <logout>
// command by sending a <hello> to a server.  Use of this element is essential
// in a connection-less environment where a server cannot return a <greeting> in
// response to a client-initiated connection.  An EPP <hello> MUST be an empty
// element with no child elements.
type Hello struct {
	// XMLName as described in https://golang.org/pkg/encoding/xml/#Unmarshal the
	// attribute is necessary to fill the namespace attribute.
	XMLName xml.Name `xml:"hello"`
}
