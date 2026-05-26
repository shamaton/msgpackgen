package msgpack

type (
	// EncResolver is a definition to resolve serialization.
	EncResolver func(i any) ([]byte, error)
	// EncToResolver resolves serialization by appending to buf.
	// If err is non-nil, encoding returns it without fallback. If handled is false,
	// encoding falls back using the original input buf. If handled is true,
	// encoding returns the resolver's buf.
	EncToResolver func(i any, buf []byte) ([]byte, bool, error)
	// DecResolver is a definition to resolve de-serialization.
	DecResolver func(data []byte, i any) (bool, error)
)

var (
	noOpEncResolver EncResolver = func(i any) ([]byte, error) {
		return nil, nil
	}
	noOpEncToResolver EncToResolver = func(i any, buf []byte) ([]byte, bool, error) {
		return buf, false, nil
	}
	noOpDecResolver DecResolver = func(data []byte, i any) (bool, error) {
		return false, nil
	}

	encAsMapResolver   = noOpEncResolver
	encAsArrayResolver = noOpEncResolver

	encAsMapToResolver   = noOpEncToResolver
	encAsArrayToResolver = noOpEncToResolver

	decAsMapResolver   = noOpDecResolver
	decAsArrayResolver = noOpDecResolver
)

// SetResolver sets generated resolvers to bridge variables.
//
// Resolver registration is intended for init/startup time. Concurrent calls to
// SetResolver or SetToResolver while Marshal/Unmarshal is running are not
// synchronized.
func SetResolver(encAsMap, encAsArray EncResolver, decAsMap, decAsArray DecResolver) {
	encAsMapResolver = encAsMap
	encAsArrayResolver = encAsArray
	encAsMapToResolver = noOpEncToResolver
	encAsArrayToResolver = noOpEncToResolver
	decAsMapResolver = decAsMap
	decAsArrayResolver = decAsArray
}

// SetToResolver sets generated resolvers that append encoded bytes to the caller's buffer.
//
// Resolver registration is intended for init/startup time. Concurrent calls to
// SetResolver or SetToResolver while Marshal/Unmarshal is running are not
// synchronized.
// A resolver that returns handled=false must leave buf unchanged.
// Passing nil for either resolver resets that side to the default no-op resolver.
func SetToResolver(encAsMap, encAsArray EncToResolver) {
	if encAsMap == nil {
		encAsMap = noOpEncToResolver
	}
	if encAsArray == nil {
		encAsArray = noOpEncToResolver
	}
	encAsMapToResolver = encAsMap
	encAsArrayToResolver = encAsArray
}
