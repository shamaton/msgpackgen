package msgpack

type (
	// EncResolver is a definition to resolve serialization.
	EncResolver func(i interface{}) ([]byte, error)
	// DecResolver is a definition to resolve de-serialization.
	DecResolver func(data []byte, i interface{}) (bool, error)
)

var (
	encAsMapResolver EncResolver = func(i interface{}) ([]byte, error) {
		return nil, nil
	}
	encAsArrayResolver = encAsMapResolver

	decAsMapResolver DecResolver = func(data []byte, i interface{}) (bool, error) {
		return false, nil
	}
	decAsArrayResolver = decAsMapResolver
)

// SetResolver sets generated resolvers to bridge variables.
func SetResolver(encAsMap, encAsArray EncResolver, decAsMap, decAsArray DecResolver) {
	encAsMapResolver = encAsMap
	encAsArrayResolver = encAsArray
	decAsMapResolver = decAsMap
	decAsArrayResolver = decAsArray
}
