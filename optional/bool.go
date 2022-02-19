package optional

import (
	"encoding/json"
	"fmt"
)

// Bool is an optional bool.
type Bool struct {
	value *bool
}

// OfBool creates an optional.Bool from a bool.
func OfBool(v bool) Bool {
	return Bool{&v}
}

// OfBoolPtr creates an optional.Bool from a bool ptr.
func OfBoolPtr(ptr *bool) Bool {
	if ptr == nil {
		return EmptyBool()
	} else {
		return OfBool(*ptr)
	}
}

// EmptyBool returns an empty optional.Bool.
func EmptyBool() Bool {
	return Bool{}
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o Bool) Get() (bool, bool) {
	if !o.Present() {
		var zero bool
		return zero, false
	}
	return *o.value, true
}

// Present returns true if there is a value wrapped by this optional.
func (o Bool) Present() bool {
	return o.value != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o Bool) If(fn func(bool)) {
	if o.Present() {
		fn(*o.value)
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o Bool) Else(v bool) bool {
	if o.Present() {
		return *o.value
	}
	return v
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o Bool) ElseZero() bool {
	var zero bool
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o Bool) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

func (o Bool) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *Bool) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value bool

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
