package error

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/constants"
	"github.com/keshu12345/overlap-avalara/pkg/customerror"
	

	//	errpkg "github.com/keshu12345/overlap-avalara/pkg/error"
	httpPkg "github.com/keshu12345/overlap-avalara/pkg/http"
	resp "github.com/keshu12345/overlap-avalara/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupGinContext returns a fresh gin.Context backed by an httptest.ResponseRecorder.
func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	return ctx, w
}

func TestNewErrorResponse_RequestInvalid(t *testing.T) {
	ctx, w := setupGinContext()

	payload := map[string]interface{}{"field": "email"}
	customErr := customerror.NewCustomErrorWithPayload(
		constants.RequestInvalid,
		"Validation failed",
		payload,
	)

	NewErrorResponse(ctx, customErr)

	assert.Equal(t, httpPkg.StatusBadRequest.Code(), w.Code)

	var body resp.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.False(t, body.IsSuccess)
	assert.Equal(t, httpPkg.StatusBadRequest.Code(), body.StatusCode)

	assert.Equal(t, "Validation failed", body.Error.Message)

	dataMap, ok := body.Error.Data.(map[string]interface{})
	require.True(t, ok, "Data should be map[string]interface{}")
	assert.Equal(t, payload["field"], dataMap["field"])

	assert.Nil(t, body.Error.Errors)
}

func TestNewErrorResponse_RequestNotValid(t *testing.T) {
	ctx, w := setupGinContext()

	customErr := customerror.NewCustomError(
		constants.RequestNotValid,
		"ShouldNotAppear",
	)

	NewErrorResponse(ctx, customErr)

	assert.Equal(t, httpPkg.StatusForbidden.Code(), w.Code)

	var body resp.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.False(t, body.IsSuccess)
	assert.Equal(t, httpPkg.StatusForbidden.Code(), body.StatusCode)

	assert.Equal(t, httpPkg.StatusForbidden.String(), body.Error.Message)

	assert.Nil(t, body.Error.Data)
	assert.Nil(t, body.Error.Errors)
}





