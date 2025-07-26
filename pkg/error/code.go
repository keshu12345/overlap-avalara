package error

import (
	"github.com/keshu12345/overlap-avalara/constants"
	"github.com/keshu12345/overlap-avalara/pkg/http"
)

var CustomCodeToHttpCodeMapping = map[constants.Code]http.StatusCode{
	RequestInvalid:  http.StatusBadRequest,
	NotFound:        http.StatusBadRequest,
	RequestNotValid: http.StatusForbidden,

	BadRequest: http.StatusBadRequest,

	ParseIntError:      http.StatusBadRequest,
	StatusUnauthorized: http.StatusUnauthorized,
}
