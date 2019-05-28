package alerting

import (
	"context"
	"testing"

	"time"

	"github.com/grafana/grafana/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

type secondCondition struct {
	firing   bool
	operator string
	matches  []*EvalMatch
	noData   bool
}

func (s *secondCondition) From() string {
	return "10m"
}

func (s *secondCondition) To() string {
	return "now-10m"
}
func (c *secondCondition) Eval(context *EvalContext) (*ConditionResult, error) {
	return &ConditionResult{Firing: c.firing, EvalMatches: c.matches, Operator: c.operator, NoDataFound: c.noData}, nil
}

func TestAlertingEvalContext(t *testing.T) {
	Convey("Eval context", t, func() {
		ctx := NewEvalContext(context.TODO(), &Rule{Conditions: []Condition{&conditionStub{firing: true}, &secondCondition{firing: true}}})

		Convey("ok -> alerting", func() {
			ctx.PrevAlertState = models.AlertStateOK
			ctx.Firing = true

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStateAlerting)
		})

		Convey("ok -> error(alerting)", func() {
			ctx.PrevAlertState = models.AlertStateOK
			ctx.Error = dummieError
			ctx.Rule.ExecutionErrorState = models.ExecutionErrorSetAlerting

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStateAlerting)
		})

		Convey("ok -> error(keep_last)", func() {
			ctx.PrevAlertState = models.AlertStateOK
			ctx.Error = dummieError
			ctx.Rule.ExecutionErrorState = models.ExecutionErrorKeepState

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStateOK)
		})

		Convey("pending -> error(keep_last)", func() {
			ctx.PrevAlertState = models.AlertStatePending
			ctx.Error = dummieError
			ctx.Rule.ExecutionErrorState = models.ExecutionErrorKeepState

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStatePending)
		})

		Convey("ok -> no_data(alerting)", func() {
			ctx.PrevAlertState = models.AlertStateOK
			ctx.Rule.NoDataState = models.NoDataSetAlerting
			ctx.NoDataFound = true

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStateAlerting)
		})

		Convey("ok -> no_data(keep_last)", func() {
			ctx.PrevAlertState = models.AlertStateOK
			ctx.Rule.NoDataState = models.NoDataKeepState
			ctx.NoDataFound = true

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStateOK)
		})

		Convey("pending -> no_data(keep_last)", func() {
			ctx.PrevAlertState = models.AlertStatePending
			ctx.Rule.NoDataState = models.NoDataKeepState
			ctx.NoDataFound = true

			ctx.Rule.State = ctx.GetNewState()
			So(ctx.Rule.State, ShouldEqual, models.AlertStatePending)
		})
		Convey("Should construct url using from, to", func() {
			ctx.StartTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
			expectedFrom := ctx.StartTime.Add(-1*(gapDuration+15*time.Minute)).UnixNano() / int64(time.Millisecond)
			expectedTo := ctx.StartTime.Add(-1*(5*time.Minute)).UnixNano() / int64(time.Millisecond)
			actualFrom, actualTo := ctx.GetFromToAsMilliseconds()
			So(actualFrom, ShouldEqual, expectedFrom)
			So(actualTo, ShouldEqual, expectedTo)
		})
	})
}
