package errors

import (
	"net/url"
)

type ResponseError struct {
	StatusCode int
	Response   *url.URL
}

type InternalError struct {
	StatusCode int
	Response   string
}

const ErrorMsg = "Internal Error. Please contact administrator for more details"
