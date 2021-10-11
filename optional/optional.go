package optional

//type Float64 struct {
//	set   atomic.Bool
//	value atomic.Float64
//}
//
//func (o *Float64) IsSet() bool {
//	return o.set.Load()
//}
//
//func (o *Float64) Value() float64 {
//	return o.value.Load()
//}
//
//func (o *Float64) Set(v float64) {
//	o.set.Store(true)
//	o.value.Store(v)
//}
//
//type Int64 struct {
//	set   atomic.Bool
//	value atomic.Int64
//}
//
//func (o *Int64) IsSet() bool {
//	return o.set.Load()
//}
//
//func (o *Int64) Value() int64 {
//	return o.value.Load()
//}
//
//func (o *Int64) Set(v int64) {
//	o.set.Store(true)
//	o.value.Store(v)
//}
//
//type String struct {
//	set   atomic.Bool
//	value atomic.String
//}
//
//func (o *String) IsSet() bool {
//	return o.set.Load()
//}
//
//func (o *String) Value() string {
//	return o.value.Load()
//}
//
//func (o *String) Set(v string) {
//	o.set.Store(true)
//	o.value.Store(v)
//}
