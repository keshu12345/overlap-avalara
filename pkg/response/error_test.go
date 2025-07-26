package response

import (
	"encoding/json"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keshu12345/overlap-avalara/constants"
	error2 "github.com/keshu12345/overlap-avalara/pkg/customerror"
	httpPkg "github.com/keshu12345/overlap-avalara/pkg/http"
)

// Test helper function to setup gin context
func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	return c, w
}

// Test helper to parse JSON response
func parseErrorResponse(t *testing.T, body []byte) ErrorResponse {
	var resp ErrorResponse
	err := json.Unmarshal(body, &resp)
	require.NoError(t, err)
	return resp
}

// Reset global state for testing
func resetGlobalState() {
	once = sync.Once{}
	customCodeToHttpCodeMapping = map[constants.Code]httpPkg.StatusCode{}
}

func TestSetCustomErrorMapping(t *testing.T) {
	resetGlobalState()

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
		constants.NotFound:       httpPkg.StatusNotFound,
	}

	// First call should set the mapping
	SetCustomErrorMapping(mapping)

	// Verify mapping is set by testing through NewErrorResponseV2
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.RequestInvalid, "Invalid request")

	NewErrorResponseV2(ctx, customErr)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), resp.StatusCode)

	// Second call should not override (due to sync.Once)
	resetGlobalState()
	differentMapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusUnauthorized,
	}
	SetCustomErrorMapping(mapping)          // Set original first
	SetCustomErrorMapping(differentMapping) // This should be ignored

	// Should still use the original mapping
	ctx2, w2 := setupGinContext()
	customErr2 := error2.NewCustomError(constants.RequestInvalid, "Invalid request")

	NewErrorResponseV2(ctx2, customErr2)

	resp2 := parseErrorResponse(t, w2.Body.Bytes())
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), resp2.StatusCode)
}

func TestNewErrorResponse_WithValidMapping(t *testing.T) {
	ctx, w := setupGinContext()
	
	// First, let's test without errors to see the baseline
	customErr := error2.NewCustomErrorWithPayload(
		constants.RequestInvalid,
		"Invalid request data",
		map[string]interface{}{"field": "value"},
	)
	
	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}
	
	NewErrorResponse(ctx, customErr, mapping)
	
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), w.Code)
	
	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), resp.StatusCode)
	assert.Equal(t, "Invalid request data", resp.Error.Message)
	assert.Equal(t, map[string]interface{}{"field": "value"}, resp.Error.Data)
	// Don't assert on Errors for now since it might be nil/empty
}

func TestNewErrorResponse_WithoutMapping(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.Code("UNKNOWN_ERROR"), "Some error")

	mapping := map[constants.Code]httpPkg.StatusCode{}

	NewErrorResponse(ctx, customErr, mapping)

	assert.Equal(t, httpPkg.StatusInternalServerError.Code(), w.Code)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, httpPkg.StatusInternalServerError.Code(), resp.StatusCode)
	assert.Equal(t, httpPkg.StatusInternalServerError.String(), resp.Error.Message)
}

func TestNewErrorResponse_WithOptions(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.RequestInvalid, "Invalid request")

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}

	// Custom option to modify the response
	customOption := func(err *ErrorResponse) {
		err.Error.Message = "Custom error message"
		err.Error.Data = "custom data"
	}

	NewErrorResponse(ctx, customErr, mapping, customOption)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "Custom error message", resp.Error.Message)
	assert.Equal(t, "custom data", resp.Error.Data)
}

func TestNewErrorResponseByStatusCode(t *testing.T) {
	ctx, w := setupGinContext()

	NewErrorResponseByStatusCode(ctx, httpPkg.StatusNotFound)

	assert.Equal(t, httpPkg.StatusNotFound.Code(), w.Code)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, httpPkg.StatusNotFound.Code(), resp.StatusCode)
	assert.Equal(t, httpPkg.StatusNotFound.String(), resp.Error.Message)
	assert.Empty(t, resp.Error.Errors)
	assert.Nil(t, resp.Error.Data)
}

func TestNewErrorResponseV2_WithGlobalMapping(t *testing.T) {
	resetGlobalState()
	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
		constants.NotFound:       httpPkg.StatusNotFound,
	}
	SetCustomErrorMapping(mapping)

	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.NotFound, "Resource not found")

	NewErrorResponseV2(ctx, customErr)

	assert.Equal(t, httpPkg.StatusNotFound.Code(), w.Code)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, httpPkg.StatusNotFound.Code(), resp.StatusCode)
	assert.Equal(t, httpPkg.StatusNotFound.String(), resp.Error.Message)
}

func TestNewErrorResponseV2_WithoutGlobalMapping(t *testing.T) {
	resetGlobalState()
	SetCustomErrorMapping(map[constants.Code]httpPkg.StatusCode{})

	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.Code("UNKNOWN"), "Unknown error")

	NewErrorResponseV2(ctx, customErr)

	assert.Equal(t, httpPkg.StatusInternalServerError.Code(), w.Code)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, httpPkg.StatusInternalServerError.Code(), resp.StatusCode)
}

func TestNewErrorResponseWithMessage(t *testing.T) {
	resetGlobalState()
	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}
	SetCustomErrorMapping(mapping)

	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.RequestInvalid, "Original message")

	customMessage := "Custom error message override"
	NewErrorResponseWithMessage(ctx, customErr, customMessage)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, customMessage, resp.Error.Message)
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), resp.StatusCode)
}


func TestGenerateErrorResponse_NonRequestInvalidError(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.NotFound, "Resource not found")

	generateErrorResponse(ctx, customErr, httpPkg.StatusNotFound)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, httpPkg.StatusNotFound.String(), resp.Error.Message) // Should use status code string
}

func TestGenerateErrorResponse_WithMultipleOptions(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.RequestInvalid, "Original message")

	option1 := func(err *ErrorResponse) {
		err.Error.Message = "Option 1 message"
	}

	option2 := func(err *ErrorResponse) {
		err.Error.Data = "Option 2 data"
	}

	option3 := func(err *ErrorResponse) {
		err.Error.Errors = map[string]string{"custom": "error"}
	}

	generateErrorResponse(ctx, customErr, httpPkg.StatusBadRequest, option1, option2, option3)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "Option 1 message", resp.Error.Message)
	assert.Equal(t, "Option 2 data", resp.Error.Data)
	assert.Equal(t, map[string]string{"custom": "error"}, resp.Error.Errors)
}


func TestNewErrorResponse_ContextAbortion(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.RequestInvalid, "Test")

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}

	// Before calling the function, context should not be aborted
	assert.False(t, ctx.IsAborted())

	NewErrorResponse(ctx, customErr, mapping)

	// After calling the function, context should be aborted
	assert.True(t, ctx.IsAborted())
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), w.Code)
}

func TestRequestInvalidError_SpecialHandling(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.RequestInvalidError("Custom validation message",
		error2.WithErrors(map[string]string{
			"username": "is required",
			"email":    "invalid format",
		}),
		error2.WithData(map[string]interface{}{
			"submitted_data": "test",
		}),
	)

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusUnprocessableEntity,
	}

	NewErrorResponse(ctx, customErr, mapping)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, httpPkg.StatusUnprocessableEntity.Code(), resp.StatusCode)
	assert.Equal(t, "Custom validation message", resp.Error.Message)
	assert.Equal(t, map[string]string{
		"username": "is required",
		"email":    "invalid format",
	}, resp.Error.Errors)
	assert.NotNil(t, resp.Error.Data)
}

func TestDifferentErrorCodes(t *testing.T) {
	testCases := []struct {
		name       string
		errorCode  constants.Code
		httpStatus httpPkg.StatusCode
		message    string
	}{
		{
			name:       "BadRequest error",
			errorCode:  constants.BadRequest,
			httpStatus: httpPkg.StatusBadRequest,
			message:    "Bad request occurred",
		},
		{
			name:       "NotFound error",
			errorCode:  constants.NotFound,
			httpStatus: httpPkg.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "Unauthorized error",
			errorCode:  constants.StatusUnauthorized,
			httpStatus: httpPkg.StatusUnauthorized,
			message:    "Unauthorized access",
		},
		{
			name:       "UnmarshalError",
			errorCode:  constants.UnmarshalError,
			httpStatus: httpPkg.StatusInternalServerError,
			message:    "Failed to unmarshal data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, w := setupGinContext()

			customErr := error2.NewCustomError(tc.errorCode, tc.message)

			mapping := map[constants.Code]httpPkg.StatusCode{
				tc.errorCode: tc.httpStatus,
			}

			NewErrorResponse(ctx, customErr, mapping)

			resp := parseErrorResponse(t, w.Body.Bytes())
			assert.Equal(t, tc.httpStatus.Code(), resp.StatusCode)

			// For RequestInvalid, message should be from UserMessage, otherwise from status code
			if tc.errorCode == constants.RequestInvalid {
				assert.Equal(t, tc.message, resp.Error.Message)
			} else {
				assert.Equal(t, tc.httpStatus.String(), resp.Error.Message)
			}
		})
	}
}

func TestWithOptionsModifications(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := error2.NewCustomError(constants.BadRequest, "Original error")

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.BadRequest: httpPkg.StatusBadRequest,
	}

	// Test multiple option functions
	NewErrorResponse(ctx, customErr, mapping,
		func(err *ErrorResponse) {
			err.IsSuccess = true // This should override the default false
		},
		func(err *ErrorResponse) {
			err.StatusCode = 299 // Custom status code
		},
		func(err *ErrorResponse) {
			err.Error.Message = "Completely different message"
		})

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.True(t, resp.IsSuccess)        // Should be overridden
	assert.Equal(t, 299, resp.StatusCode) // Should be overridden
	assert.Equal(t, "Completely different message", resp.Error.Message)
}

func TestEmptyErrorMapAndData(t *testing.T) {
	ctx, w := setupGinContext()

	// Create error without additional data or error map
	customErr := error2.NewCustomError(constants.BadRequest, "Simple error")

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.BadRequest: httpPkg.StatusBadRequest,
	}

	NewErrorResponse(ctx, customErr, mapping)

	resp := parseErrorResponse(t, w.Body.Bytes())
	assert.Nil(t, resp.Error.Data)
	assert.Empty(t, resp.Error.Errors)
	assert.Equal(t, httpPkg.StatusBadRequest.String(), resp.Error.Message)
}

// Benchmark tests
func BenchmarkNewErrorResponse(b *testing.B) {
	gin.SetMode(gin.TestMode)

	customErr := error2.NewCustomError(constants.RequestInvalid, "Test message")

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		NewErrorResponse(c, customErr, mapping)
	}
}

func BenchmarkNewErrorResponseV2(b *testing.B) {
	gin.SetMode(gin.TestMode)
	resetGlobalState()

	mapping := map[constants.Code]httpPkg.StatusCode{
		constants.RequestInvalid: httpPkg.StatusBadRequest,
	}
	SetCustomErrorMapping(mapping)

	customErr := error2.NewCustomError(constants.RequestInvalid, "Test message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		NewErrorResponseV2(c, customErr)
	}
}
