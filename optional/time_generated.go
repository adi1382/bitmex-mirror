//// Code generated by gotemplate. DO NOT EDIT.
//
package optional

//
//import (
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"time"
//)
//
//var _Time = time.Time{}
//
//// template type Optional(T)
//
//// Optional wraps a value that may or may not be nil.
//// If a value is present, it may be unwrapped to expose the underlying value.
//type Time optionalTime
//
//type optionalTime []time.Time
//
//const (
//	valueKeyTime = iota
//)
//
//// Of wraps the value in an optional.
//func OfTime(value time.Time) Time {
//	return Time{valueKeyTime: value}
//}
//
//func OfTimePtr(ptr *time.Time) Time {
//	if ptr == nil {
//		return EmptyTime()
//	} else {
//		return OfTime(*ptr)
//	}
//}
//
//// Empty returns an empty optional.
//func EmptyTime() Time {
//	return nil
//}
//
//// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
//func (o Time) Get() (value time.Time, ok bool) {
//	o.If(func(v time.Time) {
//		value = v
//		ok = true
//	})
//	return
//}
//
//// IsPresent returns true if there is a value wrapped by this optional.
//func (o Time) IsPresent() bool {
//	return o != nil
//}
//
//// If calls the function if there is a value wrapped by this optional.
//func (o Time) If(f func(value time.Time)) {
//	if o.IsPresent() {
//		f(o[valueKeyTime])
//	}
//}
//
//func (o Time) ElseFunc(f func() time.Time) (value time.Time) {
//	if o.IsPresent() {
//		o.If(func(v time.Time) { value = v })
//		return
//	} else {
//		return f()
//	}
//}
//
//// Else returns the value wrapped by this optional, or the value passed in if
//// there is no value wrapped by this optional.
//func (o Time) Else(elseValue time.Time) (value time.Time) {
//	return o.ElseFunc(func() time.Time { return elseValue })
//}
//
//// ElseZero returns the value wrapped by this optional, or the zero value of
//// the type wrapped if there is no value wrapped by this optional.
//func (o Time) ElseZero() (value time.Time) {
//	var zero time.Time
//	return o.Else(zero)
//}
//
//// String returns the string representation of the wrapped value, or the string
//// representation of the zero value of the type wrapped if there is no value
//// wrapped by this optional.
//func (o Time) String() string {
//	return fmt.Sprintf("%v", o.ElseZero())
//}
//
//// MarshalJSON marshals the value being wrapped to JSON. If there is no vale
//// being wrapped, the zero value of its type is marshaled.
//func (o Time) MarshalJSON() (data []byte, err error) {
//	return json.Marshal(o.ElseZero())
//}
//
//// UnmarshalJSON unmarshals the JSON into a value wrapped by this optional.
//func (o *Time) UnmarshalJSON(data []byte) error {
//	if string(data) == "null" {
//		*o = nil
//		return nil
//	}
//	var v time.Time
//	err := json.Unmarshal(data, &v)
//	if err != nil {
//		return err
//	}
//	*o = OfTime(v)
//	return nil
//}
//
//// MarshalXML marshals the value being wrapped to XML. If there is no vale
//// being wrapped, the zero value of its type is marshaled.
//func (o Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
//	return e.EncodeElement(o.ElseZero(), start)
//}
//
//// UnmarshalXML unmarshals the XML into a value wrapped by this optional.
//func (o *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
//	var v time.Time
//	err := d.DecodeElement(&v, &start)
//	if err != nil {
//		return err
//	}
//	*o = OfTime(v)
//	return nil
//}
