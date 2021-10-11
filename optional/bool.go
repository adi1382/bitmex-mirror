package optional

import (
	"encoding/json"
)

type Bool struct {
	set   bool
	value bool
}

func (o *Bool) IsSet() bool {
	return o.set
}

func (o *Bool) Value() bool {
	return o.value
}

func (o *Bool) Set(v bool) {
	o.set = true
	o.value = v
}

func (o Bool) MarshalJSON() (data []byte, err error) {
	if o.set {
		return json.Marshal(o.value)
	} else {
		return nil, nil
	}
}
