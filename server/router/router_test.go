package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewGinRouter_PreflightCORS(t *testing.T) {
	r, err := NewGinRouter(&config.Configuration{})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodOptions, "/any/path", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Empty(t, rec.Header().Get("Access-Control-Expose-Headers"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, rec.Header().Get("Access-Control-Max-Age"))
}

func TestNewGinRouter_NormalRequest_NoCORSOnPing(t *testing.T) {
	r, err := NewGinRouter(&config.Configuration{})
	require.NoError(t, err)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://example.com")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}
