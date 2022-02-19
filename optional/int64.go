package optional

import (
	"encoding/json"
	"fmt"
)

// Int64 is an optional int64.
type Int64 struct {
	value *int64
}

// OfInt64 creates an optional.Int64 from a int64.
func OfInt64(v int64) Int64 {
	return Int64{&v}
}

// OfInt64Ptr creates an optional.Int64 from a int64 ptr.
func OfInt64Ptr(ptr *int64) Int64 {
	if ptr == nil {
		return EmptyInt64()
	} else {
		return OfInt64(*ptr)
	}
}

// EmptyInt64 returns an empty optional.Int64.
func EmptyInt64() Int64 {
	return Int64{}
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o Int64) Get() (int64, bool) {
	if !o.Present() {
		var zero int64
		return zero, false
	}
	return *o.value, true
}

// Present returns true if there is a value wrapped by this optional.
func (o Int64) Present() bool {
	return o.value != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o Int64) If(fn func(int64)) {
	if o.Present() {
		fn(*o.value)
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o Int64) Else(v int64) int64 {
	if o.Present() {
		return *o.value
	}
	return v
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o Int64) ElseZero() int64 {
	var zero int64
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o Int64) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

func (o Int64) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *Int64) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value int64

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
