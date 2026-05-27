package msgpack

type (
	// EncResolver resolves serialization by appending to buf.
	// If err is non-nil, encoding returns it without fallback. If handled is false,
	// encoding falls back using the original input buf. If handled is true,
	// encoding returns the resolver's buf.
	EncResolver func(i any, buf []byte) ([]byte, bool, error)
	// DecResolver is a definition to resolve de-serialization.
	DecResolver func(data []byte, i any) (bool, error)
)

var (
	noOpEncResolver EncResolver = func(i any, buf []byte) ([]byte, bool, error) {
		return buf, false, nil
	}
	noOpDecResolver DecResolver = func(data []byte, i any) (bool, error) {
		return false, nil
	}

	encAsMapResolver   = noOpEncResolver
	encAsArrayResolver = noOpEncResolver

	decAsMapResolver   = noOpDecResolver
	decAsArrayResolver = noOpDecResolver
)

// SetResolver sets generated resolvers to bridge variables.
//
// Resolver registration is intended for init/startup time. Concurrent calls to
// SetResolver while Marshal/Unmarshal is running are not synchronized.
// Passing nil for either encode resolver resets that side to the default no-op
// resolver.
func SetResolver(encAsMap, encAsArray EncResolver, decAsMap, decAsArray DecResolver) {
	if encAsMap == nil {
		encAsMap = noOpEncResolver
	}
	if encAsArray == nil {
		encAsArray = noOpEncResolver
	}
	encAsMapResolver = encAsMap
	encAsArrayResolver = encAsArray
	decAsMapResolver = decAsMap
	decAsArrayResolver = decAsArray
}
