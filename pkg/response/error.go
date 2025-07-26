package response

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/constants"
	error2 "github.com/keshu12345/overlap-avalara/pkg/customerror"
	"github.com/keshu12345/overlap-avalara/pkg/http"
)

type ErrorResponse struct {
	IsSuccess  bool  `json:"is_success"`
	StatusCode int   `json:"status_code"`
	Error      Error `json:"error"`
}

type Error struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

var customCodeToHttpCodeMapping = map[constants.Code]http.StatusCode{}

var once sync.Once

// Set custom error mapping
func SetCustomErrorMapping(mapping map[constants.Code]http.StatusCode) {
	once.Do(func() {
		customCodeToHttpCodeMapping = mapping
	})
}

type Option func(*ErrorResponse)

func NewErrorResponse(
	ctx *gin.Context,
	customError error2.CustomError,
	customErrorCodeToErrorRespMapping map[constants.Code]http.StatusCode,
	options ...Option,
) {
	statusCode, ok := customErrorCodeToErrorRespMapping[customError.ErrorCode()]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	generateErrorResponse(ctx, customError, statusCode, options...)
}

func NewErrorResponseByStatusCode(ctx *gin.Context, statusCode http.StatusCode) {
	res := &ErrorResponse{
		IsSuccess:  false,
		StatusCode: statusCode.Code(),
		Error: Error{
			Message: statusCode.String(),
		},
	}

	ctx.AbortWithStatusJSON(statusCode.Code(), res)
}

func NewErrorResponseV2(
	ctx *gin.Context,
	customError error2.CustomError,
	options ...Option,
) {
	statusCode, ok := customCodeToHttpCodeMapping[customError.ErrorCode()]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	generateErrorResponse(ctx, customError, statusCode, options...)
}

func NewErrorResponseWithMessage(
	ctx *gin.Context,
	cusErr error2.CustomError,
	message string,
) {
	NewErrorResponseV2(ctx, cusErr, func(err *ErrorResponse) {
		err.Error.Message = message
	})
}

func generateErrorResponse(
	ctx *gin.Context,
	customError error2.CustomError,
	statusCode http.StatusCode,
	options ...Option,
) {

	var message string
	if customError.ErrorCode() == constants.RequestInvalid {
		message = customError.UserMessage()
	} else {
		message = statusCode.String()
	}

	res := &ErrorResponse{
		IsSuccess:  false,
		StatusCode: statusCode.Code(),
		Error: Error{
			Message: message,
			Data:    customError.ErrorData(),
			Errors:  customError.ErrorMap(),
		},
	}

	for _, option := range options {
		option(res)
	}

	ctx.AbortWithStatusJSON(statusCode.Code(), res)
}
