package intstr

import (
	"fmt"
	"strconv"
)

type IntOrString struct {
	Type   Type
	IntVal int64
	StrVal string
}

// Type represents the stored type of IntOrString.
type Type int

const (
	Int64  Type = iota // The IntOrString holds an int64.
	String             // The IntOrString holds a string.
)

func FromInt64(val int64) IntOrString {
	return IntOrString{Type: Int64, IntVal: val}
}

// FromString creates an IntOrString object with a string value.
func FromString(val string) IntOrString {
	return IntOrString{Type: String, StrVal: val}
}

// String returns the string value, or the Itoa of the int value.
func (intstr *IntOrString) String() string {
	if intstr.Type == String {
		return intstr.StrVal
	}
	return fmt.Sprintf("%d", intstr.Int64())
}

func (intstr *IntOrString) Int64() int64 {
	if intstr.Type == String {
		i, _ := strconv.ParseInt(intstr.StrVal, 10, 64)
		return i
	}
	return intstr.IntVal
}
