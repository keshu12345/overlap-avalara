package response

import (
	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/pkg/http"
)

type Success struct {
	IsSuccess  bool        `json:"is_success"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func NewSuccess(ctx *gin.Context, data interface{}) {
	res := &Success{
		IsSuccess:  true,
		StatusCode: http.StatusOK.Code(),
		Data:       data,
	}
	ctx.AbortWithStatusJSON(http.StatusOK.Code(), res)
	return
}
