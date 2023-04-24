// Package generics provide useful generic types
package generics

// Any type that is a signed integer
type SignedInt = interface {
	int | int8 | int16 | int32 | int64
}

// Any type that is an unsigned integer
type UnsignedInt = interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Any type that is a floating point number
type Float = interface {
	float32 | float64
}

// Any type that is a numeric type
type Numeric = interface {
	SignedInt | UnsignedInt | Float
}
