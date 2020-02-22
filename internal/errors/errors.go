package errors

import (
	"encoding/json"
)

type Error struct {
	StatusCode int
	Response   string
}

const (
	NilEndpointMsg    = "This endpoint does not exist on this API. Please refer to the API documentation on /api."
	InternalErrMsg    = "Internal Error. Please contact administrator for more details."
	ResponseErrMsg    = "There are no stats for the team, year combination given. Please check and make sure team and year were input correctly."
	InvalidYearErrMsg = "Please input a 4 digit year. Refer to the API documentation on /api if you need more assistance."
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
