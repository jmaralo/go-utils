package generics

type SignedInt = interface {
	int | int8 | int16 | int32 | int64
}

type UnsignedInt = interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float = interface {
	float32 | float64
}

type Numeric = interface {
	SignedInt | UnsignedInt | Float
}
