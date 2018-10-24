package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting/riot"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/leibowitz/moment"
	"math"
)

var (
	eslog log.Logger
)

type ElasticsearchExecutor struct {
	*models.DataSource
}

const MAX_ES_CLUSTER_LIMIT = 8

func NewElasticsearchExecutor(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	return &ElasticsearchExecutor{}, nil
}

func init() {
	eslog = log.New("tsdb.elasticsearch")
	tsdb.RegisterTsdbQueryEndpoint("elasticsearch", NewElasticsearchExecutor)
}

func getIndex(pattern string, interval string, timeRange *tsdb.TimeRange) string {
	if interval == "" {
		return pattern
	}

	indexes := []string{}
	indexParts := strings.Split(strings.TrimLeft(pattern, "["), "]")
	indexBase := indexParts[0]
	if len(indexParts) <= 1 {
		return pattern
	}
	indexDateFormat := indexParts[1]

	if interval == "" {
		return pattern
	}

	type TimeOperation struct {
		StartOf   func(*moment.Moment) *moment.Moment
		Operation func(*moment.Moment) *moment.Moment
	}

	operations := map[string]TimeOperation{
		"Hourly": {
			StartOf:   func(m *moment.Moment) *moment.Moment { return m.StartOf("hour") },
			Operation: func(m *moment.Moment) *moment.Moment { return m.AddHours(1) },
		},
		"Daily": {
			StartOf:   func(m *moment.Moment) *moment.Moment { return m.StartOfDay() },
			Operation: func(m *moment.Moment) *moment.Moment { return m.AddDay() },
		},
		"Weekly": {
			StartOf:   func(m *moment.Moment) *moment.Moment { return m.StartOfWeek() },
			Operation: func(m *moment.Moment) *moment.Moment { return m.AddWeeks(1) },
		},
		"Monthly": {
			StartOf:   func(m *moment.Moment) *moment.Moment { return m.StartOfMonth() },
			Operation: func(m *moment.Moment) *moment.Moment { return m.AddMonths(1) },
		},
		"Yearly": {
			StartOf:   func(m *moment.Moment) *moment.Moment { return m.StartOfYear() },
			Operation: func(m *moment.Moment) *moment.Moment { return m.AddYears(1) },
		},
	}

	if _, ok := operations[interval]; !ok {
		eslog.Error("Unknown ElasticSearch Interval: %s", interval)
		return pattern
	}

	start := operations[interval].StartOf(moment.NewMoment(timeRange.MustGetFrom()).UTC())
	end := operations[interval].StartOf(moment.NewMoment(timeRange.MustGetTo()).UTC())

	indexes = append(indexes, fmt.Sprintf("%s%s", indexBase, start.UTC().Format(indexDateFormat)))
	for start.IsBefore(*end) {
		start = operations[interval].Operation(start)
		indexes = append(indexes, fmt.Sprintf("%s%s", indexBase, start.UTC().Format(indexDateFormat)))
	}
	limitedIndexes := indexes[int(math.Max(0, float64(len(indexes) - MAX_ES_CLUSTER_LIMIT))):]
	return strings.Join(limitedIndexes, ",")
}

func (e *ElasticsearchExecutor) buildRequest(queryInfo *tsdb.Query, timeRange *tsdb.TimeRange) (*http.Request, error) {
	indexIntervalJson := queryInfo.DataSource.JsonData.Get("interval")
	var indexInterval = ""
	if indexIntervalJson.Interface() != nil {
		indexIntervalStr, err := indexIntervalJson.String()
		if err != nil {
			return nil, err
		}
		indexInterval = indexIntervalStr
	}

	index := getIndex(queryInfo.DataSource.Database, indexInterval, timeRange)

	esRequestURL := fmt.Sprintf("%s/%s/_search", queryInfo.DataSource.Url, index)

	if queryInfo.Model == nil {
		return nil, fmt.Errorf("Invalid (nil) ES Request Model Provided!")
	}

	esRequestModel := &RequestModel{}
	rawModel, err := queryInfo.Model.MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawModel, esRequestModel)
	if err != nil {
		return nil, err
	}

	esRequestModel.ESVersion, err = queryInfo.DataSource.JsonData.Get("esVersion").Float64()
	if err != nil {
		return nil, err
	}

	esRequestJSON, err := esRequestModel.buildQueryJSON(timeRange)
	if err != nil {
		return nil, err
	}

	intervalCalculator := tsdb.NewIntervalCalculator(&tsdb.IntervalOptions{})
	interval := intervalCalculator.Calculate(timeRange, time.Millisecond)

	esRequestJSON = strings.Replace(esRequestJSON, "$interval", interval.Text, 1)
	esRequestJSON = strings.Replace(esRequestJSON, "$__interval_ms", strconv.FormatInt(interval.Value.Nanoseconds()/int64(time.Millisecond), 10), 1)
	esRequestJSON = strings.Replace(esRequestJSON, "$__interval", interval.Text, 1)

	reader := strings.NewReader(esRequestJSON)
	req, err := http.NewRequest("GET", esRequestURL, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if queryInfo.DataSource.BasicAuth {
		req.SetBasicAuth(queryInfo.DataSource.BasicAuthUser, queryInfo.DataSource.BasicAuthPassword)
	}

	return req, nil
}

func (e *ElasticsearchExecutor) Query(ctx context.Context, dsInfo *models.DataSource, queryContext *tsdb.TsdbQuery) (*tsdb.Response, error) {
	result := &tsdb.Response{}
	result.Results = make(map[string]*tsdb.QueryResult)
	rtpOk, err := riot.IsRTPHealthy(dsInfo.Url)
	if err != nil {
		eslog.Warn("Error Getting RTP Circuit Breaker Status: %s\n", err.Error())
	}

	if !rtpOk {
		return nil, fmt.Errorf("RTP Is Currently Unhealthy, Raising an Execution Error.")
	}

	for _, q := range queryContext.Queries {
		if q.DataSource == nil {
			return nil, fmt.Errorf("Invalid (nil) DataSource Provided")
		}

		if q.DataSource.JsonData == nil {
			return nil, fmt.Errorf("Invalid (nil) JsonData Provided")
		}

		esRequest, err := e.buildRequest(q, queryContext.TimeRange)
		if err != nil {
			return nil, err
		}

		httpClient, err := dsInfo.GetHttpClient()
		if err != nil {
			return nil, err
		}

		resp, err := httpClient.Do(esRequest)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		rBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body (%s): %s", strconv.Quote(string(rBody)), err)
		}

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Failed to get metrics: %s (%s)", strconv.Quote(string(rBody)), resp.Status)
		}

		result.Results[q.RefId], err = parseQueryResult(rBody, getPreferredNamesForQueries(q), getFilteredMetrics(q))
		if err != nil {
			return result, err
		}

		result.Results[q.RefId].RefId = q.RefId
	}

	return result, nil
}

func getPreferredNamesForQueries(query *tsdb.Query) NameMap {
	names := NameMap{}
	alias := query.Model.Get("alias").MustString("")

	esRequestModel := &RequestModel{}
	rawModel, err := query.Model.MarshalJSON()
	if err != nil {
		return names
	}

	err = json.Unmarshal(rawModel, esRequestModel)
	if err != nil {
		return names
	}

	for _, metric := range esRequestModel.Metrics {
		if alias != "" {
			names[metric.ID] = Name{
				Value: alias,
			}
		} else {
			names[metric.ID] = metric.GetName()
		}
	}

	return names
}

func getFilteredMetrics(query *tsdb.Query) FilterMap {
	filters := FilterMap{}

	esRequestModel := &RequestModel{}
	rawModel, err := query.Model.MarshalJSON()
	if err != nil {
		return filters
	}

	err = json.Unmarshal(rawModel, esRequestModel)
	if err != nil {
		return filters
	}

	for _, metric := range esRequestModel.Metrics {
		filters[metric.ID] = metric.Hide
	}

	return filters
}
