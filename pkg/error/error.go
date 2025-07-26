package error

import (
	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/constants"
	"github.com/keshu12345/overlap-avalara/pkg/customerror"
	"github.com/keshu12345/overlap-avalara/pkg/response"
)

const (
	BadRequest          constants.Code = "BAD_REQUEST"
	NotFound            constants.Code = "NOT_FOUND"
	RequestNotValid     constants.Code = "REQUEST_NOT_VALID"
	RequestInvalid      constants.Code = "REQUEST_INVALID"
	UnmarshalError      constants.Code = "UNMARSHAl_ERROR"
	MarshalError        constants.Code = "MARSHAL_ERR"
	ParseIntError       constants.Code = "PARSE_INT_ERROR"
	DataNotFoundDbError constants.Code = "DATA_NOT_FOUND_DB_ERROR"
	GoroutineError      constants.Code = "GOROUTINE_ERROR"
	ParseFilesError     constants.Code = "PARSE_FILES_ERROR"
	NotFoundMapError    constants.Code = "NOT_FOUND_MAP_ERROR"
	UrlError            constants.Code = "URL_ERROR"
	StatusUnauthorized  constants.Code = "UNAUTHORIZED_ERROR"
)

func NewErrorResponse(ctx *gin.Context, cusErr customerror.CustomError) {
	response.NewErrorResponse(ctx, cusErr, CustomCodeToHttpCodeMapping)
}
