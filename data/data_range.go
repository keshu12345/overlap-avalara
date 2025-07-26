package data

import "time"

type OverlapRequest struct {
	Range1 DateRange `json:"range1" binding:"required"`
	Range2 DateRange `json:"range2" binding:"required"`
}

type DateRange struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
}
