package optional

import (
	"encoding/json"
	"fmt"
)

// String is an optional string.
type String struct {
	value *string
}

// OfString creates an optional.String from a string.
func OfString(v string) String {
	return String{&v}
}

// OfStringPtr creates an optional.String from a string ptr.
func OfStringPtr(ptr *string) String {
	if ptr == nil {
		return EmptyString()
	} else {
		return OfString(*ptr)
	}
}

// EmptyString returns an empty optional.String.
func EmptyString() String {
	return String{}
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o String) Get() (string, bool) {
	if !o.Present() {
		var zero string
		return zero, false
	}
	return *o.value, true
}

// Present returns true if there is a value wrapped by this optional.
func (o String) Present() bool {
	return o.value != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o String) If(fn func(string)) {
	if o.Present() {
		fn(*o.value)
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o String) Else(v string) string {
	if o.Present() {
		return *o.value
	}
	return v
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o String) ElseZero() string {
	var zero string
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o String) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

func (o String) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *String) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
