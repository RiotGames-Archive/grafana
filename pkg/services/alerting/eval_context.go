package alerting

import (
	"context"
	"fmt"
	"time"

	"strings"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type EvalContext struct {
	Firing         bool
	IsTestRun      bool
	EvalMatches    []*EvalMatch
	Logs           []*ResultLogEntry
	Error          error
	ConditionEvals string
	StartTime      time.Time
	EndTime        time.Time
	Rule           *Rule
	log            log.Logger

	dashboardRef *m.DashboardRef

	ImagePublicUrl  string
	ImageOnDiskPath string
	NoDataFound     bool
	PrevAlertState  m.AlertStateType

	Ctx context.Context
}

func NewEvalContext(alertCtx context.Context, rule *Rule) *EvalContext {
	return &EvalContext{
		Ctx:            alertCtx,
		StartTime:      time.Now(),
		Rule:           rule,
		Logs:           make([]*ResultLogEntry, 0),
		EvalMatches:    make([]*EvalMatch, 0),
		log:            log.New("alerting.evalContext"),
		PrevAlertState: rule.State,
	}
}

type StateDescription struct {
	Color string
	Text  string
	Data  string
}

func (c *EvalContext) GetStateModel() *StateDescription {
	switch c.Rule.State {
	case m.AlertStateOK:
		return &StateDescription{
			Color: "#36a64f",
			Text:  "OK",
		}
	case m.AlertStateNoData:
		return &StateDescription{
			Color: "#888888",
			Text:  "No Data",
		}
	case m.AlertStateAlerting:
		return &StateDescription{
			Color: "#D63232",
			Text:  "Alerting",
		}
	default:
		panic("Unknown rule state " + c.Rule.State)
	}
}

func (c *EvalContext) ShouldUpdateAlertState() bool {
	return c.Rule.State != c.PrevAlertState
}

func (a *EvalContext) GetDurationMs() float64 {
	return float64(a.EndTime.Nanosecond()-a.StartTime.Nanosecond()) / float64(1000000)
}

func (c *EvalContext) GetNotificationTitle() string {
	return "[" + c.GetStateModel().Text + "] " + c.Rule.Name
}

func (c *EvalContext) GetDashboardUID() (*m.DashboardRef, error) {
	if c.dashboardRef != nil {
		return c.dashboardRef, nil
	}

	uidQuery := &m.GetDashboardRefByIdQuery{Id: c.Rule.DashboardId}
	if err := bus.Dispatch(uidQuery); err != nil {
		return nil, err
	}

	c.dashboardRef = uidQuery.Result
	return c.dashboardRef, nil
}

const urlFormat = "%s?fullscreen=true&edit=true&tab=alert&panelId=%d&orgId=%d&from=%d&to=%d"

func (c *EvalContext) GetRuleUrl() (string, error) {
	if c.IsTestRun {
		return setting.AppUrl, nil
	}

	ref, err := c.GetDashboardUID()
	if err != nil {
		return "", err
	}
	from, to := c.GetFromToAsMilliseconds()
	return fmt.Sprintf(urlFormat, m.GetFullDashboardUrl(ref.Uid, ref.Slug), c.Rule.PanelId, c.Rule.OrgId, from, to), nil
}

const gapDuration = time.Minute * 30
const maxSlide = time.Hour*24*2 + gapDuration // Two days + gap duration
const maxDuration time.Duration = 1<<63 - 1

func (c *EvalContext) GetFromToAsMilliseconds() (int64, int64) {
	toSlide := maxDuration
	window := 0 * time.Minute
	var from, to time.Time

	for _, cond := range c.Rule.Conditions {
		queryWindow, err := time.ParseDuration(cond.From())
		if err == nil && queryWindow > window {
			window = queryWindow
		}
		if strings.Contains(cond.To(), "now-") {
			end := strings.Replace(cond.To(), "now-", "", -1)
			endSlide, err := time.ParseDuration(end)
			if err != nil {
				endSlide = maxDuration
			}
			if endSlide < toSlide {
				toSlide = endSlide
			}
		}
	}

	if toSlide == maxDuration {
		toSlide = 0
	}

	totalSlide := window + toSlide + gapDuration
	if totalSlide > maxSlide {
		totalSlide = maxSlide
	}

	from = c.StartTime.Add(-1 * totalSlide)
	to = c.StartTime.Add(-1 * toSlide)
	toMillis := func(aTime time.Time) int64 {
		return aTime.UnixNano() / int64(time.Millisecond)
	}
	return toMillis(from), toMillis(to)
}

func (c *EvalContext) GetNewState() m.AlertStateType {
	if c.Error != nil {
		c.log.Error("Alert Rule Result Error",
			"ruleId", c.Rule.Id,
			"name", c.Rule.Name,
			"error", c.Error,
			"changing state to", c.Rule.ExecutionErrorState.ToAlertState())

		if c.Rule.ExecutionErrorState == m.ExecutionErrorKeepState {
			return c.PrevAlertState
		}
		return c.Rule.ExecutionErrorState.ToAlertState()

	} else if c.Firing {
		return m.AlertStateAlerting

	} else if c.NoDataFound {
		c.log.Info("Alert Rule returned no data",
			"ruleId", c.Rule.Id,
			"name", c.Rule.Name,
			"changing state to", c.Rule.NoDataState.ToAlertState())

		if c.Rule.NoDataState == m.NoDataKeepState {
			return c.PrevAlertState
		}
		return c.Rule.NoDataState.ToAlertState()
	}

	return m.AlertStateOK
}
