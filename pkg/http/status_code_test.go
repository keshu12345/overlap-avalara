package http
// pkg/http/http_test.go

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCode_Code(t *testing.T) {
	assert.Equal(t, 200, StatusOK.Code())
	assert.Equal(t, 201, StatusCreated.Code())
	assert.Equal(t, 404, StatusNotFound.Code())
}

func TestStatusCode_String_CustomMapping(t *testing.T) {
	// Entries present in StatusCodeToStringMap
	assert.Equal(t, "Success", StatusOK.String())
	assert.Equal(t, "Invalid Request", StatusBadRequest.String())
	assert.Equal(t, "Something went wrong", StatusInternalServerError.String())
	assert.Equal(t, "Validation failed", StatusForbidden.String())
	assert.Equal(t, "No Content", StatusNoContent.String())
}

func TestStatusCode_String_FallbackToHTTP(t *testing.T) {
	// 201 Created is not in the custom map, so should use net/http.StatusText
	expected := http.StatusText(http.StatusCreated)
	assert.Equal(t, expected, StatusCreated.String())

	// Pick an arbitrary code not in map, e.g. 418 ("I'm a teapot")
	teapot := StatusCode(418)
	assert.Equal(t, http.StatusText(418), teapot.String())
}

func TestStatusCode_Is2xx(t *testing.T) {
	assert.True(t, StatusOK.Is2xx())
	assert.True(t, StatusCreated.Is2xx())

	assert.False(t, StatusCode(199).Is2xx())
	assert.False(t, StatusCode(300).Is2xx())
}

func TestStatusCode_Is3xx(t *testing.T) {
	assert.True(t, StatusMovedPermanently.Is3xx())
	assert.True(t, StatusFound.Is3xx())

	assert.False(t, StatusCode(299).Is3xx())
	assert.False(t, StatusCode(400).Is3xx())
}

func TestStatusCode_Is4xx(t *testing.T) {
	assert.True(t, StatusBadRequest.Is4xx())
	assert.True(t, StatusTooManyRequests.Is4xx())

	assert.False(t, StatusCode(399).Is4xx())
	assert.False(t, StatusCode(500).Is4xx())
}

func TestStatusCode_Is5xx(t *testing.T) {
	assert.True(t, StatusInternalServerError.Is5xx())
	assert.True(t, StatusBadGateway.Is5xx())

	assert.False(t, StatusCode(499).Is5xx())
	assert.False(t, StatusCode(600).Is5xx())
}

func TestAPIMethod_String(t *testing.T) {
	tests := []struct {
		method APIMethod
		want   string
	}{
		{APIGet, "GET"},
		{APIPost, "POST"},
		{APIPut, "PUT"},
		{APIDelete, "DELETE"},
		{APIPatch, "PATCH"},
		{APIHead, "HEAD"},
		{APIOptions, "OPTIONS"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.method.String())
	}
}
