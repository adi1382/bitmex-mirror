// Code generated by gotemplate. DO NOT EDIT.

package optional

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

var _Int = time.Time{}

// template type Optional(T)

// Optional wraps a value that may or may not be nil.
// If a value is present, it may be unwrapped to expose the underlying value.
type Int optionalInt

type optionalInt []int

const (
	valueKeyInt = iota
)

// Of wraps the value in an optional.
func OfInt(value int) Int {
	return Int{valueKeyInt: value}
}

func OfIntPtr(ptr *int) Int {
	if ptr == nil {
		return EmptyInt()
	} else {
		return OfInt(*ptr)
	}
}

// Empty returns an empty optional.
func EmptyInt() Int {
	return nil
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o Int) Get() (value int, ok bool) {
	o.If(func(v int) {
		value = v
		ok = true
	})
	return
}

// IsPresent returns true if there is a value wrapped by this optional.
func (o Int) IsPresent() bool {
	return o != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o Int) If(f func(value int)) {
	if o.IsPresent() {
		f(o[valueKeyInt])
	}
}

func (o Int) ElseFunc(f func() int) (value int) {
	if o.IsPresent() {
		o.If(func(v int) { value = v })
		return
	} else {
		return f()
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o Int) Else(elseValue int) (value int) {
	return o.ElseFunc(func() int { return elseValue })
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o Int) ElseZero() (value int) {
	var zero int
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o Int) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

// MarshalJSON marshals the value being wrapped to JSON. If there is no vale
// being wrapped, the zero value of its type is marshaled.
func (o Int) MarshalJSON() (data []byte, err error) {
	return json.Marshal(o.ElseZero())
}

// UnmarshalJSON unmarshals the JSON into a value wrapped by this optional.
func (o *Int) UnmarshalJSON(data []byte) error {
	var v int
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = OfInt(v)
	return nil
}

// MarshalXML marshals the value being wrapped to XML. If there is no vale
// being wrapped, the zero value of its type is marshaled.
func (o Int) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(o.ElseZero(), start)
}

// UnmarshalXML unmarshals the XML into a value wrapped by this optional.
func (o *Int) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	*o = OfInt(v)
	return nil
}
