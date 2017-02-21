package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
)

type GMetricsdExecutor struct {
	*models.DataSource
	HttpClient *http.Client
}

func NewGMetricsdExecutor(dsInfo *models.DataSource) (tsdb.Executor, error) {
	client, err := dsInfo.GetHttpClient()
	if err != nil {
		return nil, err
	}

	return &GMetricsdExecutor{
		DataSource: dsInfo,
		HttpClient: client,
	}, nil
}

func (e *GMetricsdExecutor) buildRequest(queryInfo *tsdb.Query, timeRange *tsdb.TimeRange) (*http.Request, error) {

	if queryInfo.Model == nil {
		return nil, fmt.Errorf("Invalid (nil) GMetricsd Request Model Provided!")
	}

	requestModel := &GMetricsdRequestModel{}
	rawModel, err := queryInfo.Model.MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawModel, requestModel)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(queryInfo.DataSource.Url)
	if err != nil {
		return nil, err
	}

	if queryInfo.DataSource.BasicAuth {
		parsedURL.User = url.UserPassword(queryInfo.DataSource.BasicAuthUser, queryInfo.DataSource.BasicAuthPassword)
	}

	requestURL := fmt.Sprintf("%s/v1/api/metric/history?uuid=%s&start=%d&end=%d&numslices=100", parsedURL.String(), requestModel.Metric.UUID, timeRange.GetFromAsMsEpoch()/1000, timeRange.GetToAsMsEpoch()/1000)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (e *GMetricsdExecutor) Execute(ctx context.Context, queries tsdb.QuerySlice, context *tsdb.QueryContext) *tsdb.BatchResult {
	result := &tsdb.BatchResult{}
	result.QueryResults = make(map[string]*tsdb.QueryResult)

	if context == nil {
		return result.WithError(fmt.Errorf("Nil Context provided to GMetricsdExecutor"))
	}

	for _, q := range context.Queries {
		if q.DataSource == nil {
			return result.WithError(fmt.Errorf("Invalid (nil) DataSource Provided"))
		}

		request, err := e.buildRequest(q, context.TimeRange)
		if err != nil {
			return result.WithError(err)
		}

		resp, err := e.HttpClient.Do(request)
		if err != nil {
			return result.WithError(err)
		}
		defer resp.Body.Close()

		rBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result.WithError(fmt.Errorf("Failed to read response body (%s): %s", strconv.Quote(string(rBody)), err))
		}

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return result.WithError(fmt.Errorf("Failed to get metrics: %s (%s)", strconv.Quote(string(rBody)), resp.Status))
		}

		gmetricsdHistory := &GMetricsdHistoryResponse{}
		err = json.Unmarshal(rBody, gmetricsdHistory)

		parsedQueryResult, err := gmetricsdHistory.ToTsdbQueryResult()
		if err != nil {
			return result.WithError(fmt.Errorf("Failed to parse gmetricsd history response: %s", err.Error()))
		}
		parsedQueryResult.RefId = q.RefId
		result.QueryResults[q.RefId] = parsedQueryResult
	}

	return result
}
