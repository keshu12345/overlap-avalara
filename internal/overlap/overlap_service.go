package overlap

import (
	"github.com/keshu12345/overlap-avalara/data"
	"github.com/keshu12345/overlap-avalara/logger"
)

// mockery --exported --name=OverlapService --case underscore --output ../../mocks/overlapservice
type OverlapService interface {
	Check(r1, r2 data.DateRange) bool
}

type overlapService struct {
	Logger logger.Logger
}

func New(logger logger.Logger) OverlapService {
	return &overlapService{
		Logger: logger,
	}
}

func (os *overlapService) Check(r1, r2 data.DateRange) bool {
	os.Logger.Info("Checking time range  with overlapservice")
	return r1.Start.Before(r2.End) && r2.Start.Before(r1.End)
}
