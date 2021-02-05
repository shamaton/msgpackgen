package msgpack

type (
	EncResolver func(i interface{}) ([]byte, error)
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

func SetResolver(encAsMap, encAsArray EncResolver, decAsMap, decAsArray DecResolver) {
	encAsMapResolver = encAsMap
	encAsArrayResolver = encAsArray
	decAsMapResolver = decAsMap
	decAsArrayResolver = decAsArray
}
