package optional

import (
	"encoding/json"
	"fmt"
	"time"
)

// Time is an optional time.Time.Time.
type Time struct {
	value *time.Time
}

// OfTime creates an optional.Time from a time.Time.
func OfTime(v time.Time) Time {
	return Time{&v}
}

// OfTimePtr creates an optional.Time from a time.Time ptr.
func OfTimePtr(ptr *time.Time) Time {
	if ptr == nil {
		return EmptyTime()
	} else {
		return OfTime(*ptr)
	}
}

// EmptyTime returns an empty optional.Time.
func EmptyTime() Time {
	return Time{}
}

// Get returns the value wrapped by this optional, and an ok signal for whether a value was wrapped.
func (o Time) Get() (time.Time, bool) {
	if !o.Present() {
		var zero time.Time
		return zero, false
	}
	return *o.value, true
}

// Present returns true if there is a value wrapped by this optional.
func (o Time) Present() bool {
	return o.value != nil
}

// If calls the function if there is a value wrapped by this optional.
func (o Time) If(fn func(time.Time)) {
	if o.Present() {
		fn(*o.value)
	}
}

// Else returns the value wrapped by this optional, or the value passed in if
// there is no value wrapped by this optional.
func (o Time) Else(v time.Time) time.Time {
	if o.Present() {
		return *o.value
	}
	return v
}

// ElseZero returns the value wrapped by this optional, or the zero value of
// the type wrapped if there is no value wrapped by this optional.
func (o Time) ElseZero() time.Time {
	var zero time.Time
	return o.Else(zero)
}

// String returns the string representation of the wrapped value, or the string
// representation of the zero value of the type wrapped if there is no value
// wrapped by this optional.
func (o Time) String() string {
	return fmt.Sprintf("%v", o.ElseZero())
}

func (o Time) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *Time) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value time.Time

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
