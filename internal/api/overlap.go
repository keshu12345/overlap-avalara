package api

import (
	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/data"
	"github.com/keshu12345/overlap-avalara/pkg/customerror"
	"github.com/keshu12345/overlap-avalara/pkg/error"
	"github.com/keshu12345/overlap-avalara/pkg/response"
)

func CheckOverlap(c *gin.Context) {
	var req data.OverlapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cusErr := customerror.NewCustomError(error.BadRequest, err.Error())
		appLogger.Errorf("Unable to bind with json body :%v", cusErr)
		error.NewErrorResponse(c, cusErr)
		return
	}

	isOverlap := overlapService.Check(req.Range1, req.Range2)
	appLogger.Infof("isOverlap the time range %v", isOverlap)
	response.NewSuccess(c, isOverlap)
}
