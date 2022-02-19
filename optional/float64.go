package optional

import (
	"encoding/json"
	"fmt"
)

// Float64 is an optional float64.
type Float64 struct {
	value *float64
}

// OfFloat64 creates an optional.Float64 from a float64.
func OfFloat64(v float64) Float64 {
	return Float64{&v}
}

// OfFloat64Ptr creates an optional.Float64 from a float64 ptr.
func OfFloat64Ptr(ptr *float64) Float64 {
	if ptr == nil {
		return EmptyFloat64()
	} else {
		return OfFloat64(*ptr)
	}
}

// EmptyFloat64 returns an empty optional.Float64.
func EmptyFloat64() Float64 {
	return Float64{}
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o Float64) Get() (float64, bool) {
	if !o.Present() {
		var zero float64
		return zero, false
	}
	return *o.value, true
}

// Present returns true if there is a value wrapped by this optional.
func (o Float64) Present() bool {
	return o.value != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o Float64) If(fn func(float64)) {
	if o.Present() {
		fn(*o.value)
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o Float64) Else(v float64) float64 {
	if o.Present() {
		return *o.value
	}
	return v
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o Float64) ElseZero() float64 {
	var zero float64
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o Float64) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

func (o Float64) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *Float64) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value float64

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
