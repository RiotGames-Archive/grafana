package riot

import (
	"fmt"
	"strconv"

	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/tsdb"
)

type MetricObject struct {
	UUID string
}

type GMetricsdRequestModel struct {
	Metric MetricObject `json:"metricObject"`
}

type DataPointValue struct {
	Timestamp string
	Value     interface{}
}

type DataPoint struct {
	StartTime          string
	EndTime            string
	Min                DataPointValue
	Max                DataPointValue
	Average            interface{}
	Midrange           interface{}
	NumberOfDataPoints int
}

type HistoryObject struct {
	UUID       string `json:"Uuid"`
	TimeSeries map[string]DataPoint
	PrettyName string
	Name       string
	Host       Host
}

type Host struct {
	Hostname string
	Hostid   string
}

type GMetricsdHistoryResponse struct {
	History HistoryObject
}

func (history *GMetricsdHistoryResponse) ToTsdbQueryResult() (*tsdb.QueryResult, error) {
	queryRes := tsdb.NewQueryResult()

	var points tsdb.TimeSeriesPoints
	for ts, data := range history.History.TimeSeries {
		var valueRow [2]null.Float
		fValue, ok := data.Average.(float64)
		if !ok {
			continue
		}
		valueRow[0] = parseValue(fValue)

		tsFloat, err := strconv.ParseFloat(ts, 64)
		if err != nil {
			riotlog.Error(fmt.Sprintf("Error parsing timestamp in gmetricsd alert handler: %s", err.Error()))
			continue
		}
		valueRow[1] = parseValue(tsFloat)

		points = append(points, valueRow)
	}

	ts := &tsdb.TimeSeries{
		Name:   fmt.Sprintf("%s [%s]", history.History.PrettyName, history.History.Host.Hostname),
		Points: points,
	}

	queryRes.Series = append(queryRes.Series, ts)
	return queryRes, nil
}

func parseValue(value float64) null.Float {
	return null.FloatFrom(float64(value))
}
