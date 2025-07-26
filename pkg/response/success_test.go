// pkg/response/success_test.go
package response

import (
	"encoding/json"
	"testing"

	httpPkg "github.com/keshu12345/overlap-avalara/pkg/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSuccess_NoData(t *testing.T) {
	ctx, w := setupGinContext()

	NewSuccess(ctx, nil)

	assert.Equal(t, httpPkg.StatusOK.Code(), w.Code)

	var resp Success
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.True(t, resp.IsSuccess, "IsSuccess should be true")
	assert.Equal(t, httpPkg.StatusOK.Code(), resp.StatusCode, "StatusCode mismatch")
	assert.Nil(t, resp.Data, "Data should be nil when passed nil")
}

func TestNewSuccess_WithData(t *testing.T) {
	ctx, w := setupGinContext()

	original := map[string]interface{}{
		"foo": "bar",
		"num": 42,
	}

	NewSuccess(ctx, original)

	assert.Equal(t, httpPkg.StatusOK.Code(), w.Code)

	var resp Success
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.True(t, resp.IsSuccess, "IsSuccess should be true")
	assert.Equal(t, httpPkg.StatusOK.Code(), resp.StatusCode, "StatusCode mismatch")

	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok, "Data should be a map[string]interface{}")

	assert.Equal(t, "bar", dataMap["foo"], "foo value mismatch")
	assert.EqualValues(t, 42, dataMap["num"], "num value mismatch")
}
