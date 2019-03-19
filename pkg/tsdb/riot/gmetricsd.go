package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/context/ctxhttp"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
)

type GMetricsdExecutor struct {
}

func NewGMetricsdExecutor(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	return &GMetricsdExecutor{}, nil
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

func (e *GMetricsdExecutor) Query(ctx context.Context, dsInfo *models.DataSource, queryContext *tsdb.TsdbQuery) (*tsdb.Response, error) {
	result := &tsdb.Response{}
	result.Results = make(map[string]*tsdb.QueryResult)

	for _, q := range queryContext.Queries {
		if q.DataSource == nil {
			return nil, fmt.Errorf("Invalid (nil) DataSource Provided")
		}

		request, err := e.buildRequest(q, queryContext.TimeRange)
		if err != nil {
			return nil, err
		}

		httpClient, err := dsInfo.GetHttpClient()
		if err != nil {
			return nil, err
		}

		resp, err := ctxhttp.Do(ctx, httpClient, request)
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

		gmetricsdHistory := &GMetricsdHistoryResponse{}
		err = json.Unmarshal(rBody, gmetricsdHistory)

		parsedQueryResult, err := gmetricsdHistory.ToTsdbQueryResult()
		if err != nil {
			return nil, fmt.Errorf("Failed to parse gmetricsd history response: %s", err.Error())
		}
		parsedQueryResult.RefId = q.RefId
		result.Results[q.RefId] = parsedQueryResult
	}

	return result, nil
}
