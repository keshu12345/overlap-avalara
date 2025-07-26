package api

import (
	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/internal/overlap"
	"github.com/keshu12345/overlap-avalara/logger"
)

var overlapService overlap.OverlapService

var appLogger logger.Logger

func RegisterEndpoint(g *gin.Engine, os overlap.OverlapService, logger logger.Logger) {

	overlapService = os
	appLogger = logger

	v1 := g.Group("/api/v1")
	{

		v1.POST("/overlap-check", CheckOverlap)
	}
}
