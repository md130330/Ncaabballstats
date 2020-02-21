package errors

import (
	"encoding/json"
)

type Error struct {
	StatusCode int
	Response   string
}

const (
	NilEndpointMsg    = "This endpoint does not exist on this API. Please refer to the API documentation."
	InternalErrMsg    = "Internal Error. Please contact administrator for more details."
	ResponseErrMsg    = "There is an error in the team or year. Please double check and try again."
	InvalidYearErrMsg = "Please input a 4 digit year. Refer to the API documentation if you need more assistance."
)

func NilEndpoint() []byte {
	nilEndpoint := Error{404, NilEndpointMsg}
	errJson, _ := json.Marshal(nilEndpoint)
	return errJson
}

func InternalErr() []byte {
	internalErr := Error{500, InternalErrMsg}
	errJson, _ := json.Marshal(internalErr)
	return errJson
}

func ResponseErr() []byte {
	responseErr := Error{404, ResponseErrMsg}
	errJson, _ := json.Marshal(responseErr)
	return errJson
}

func InvalidYearErr() []byte {
	invalidYearErr := Error{400, InvalidYearErrMsg}
	errJson, _ := json.Marshal(invalidYearErr)
	return errJson
}
