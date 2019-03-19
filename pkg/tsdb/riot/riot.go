package riot

import (
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/grafana/grafana/pkg/tsdb/elasticsearch"
)

var (
	riotlog log.Logger
)

func init() {
	riotlog = log.New("tsdb.riot")

	tsdb.RegisterTsdbQueryEndpoint("riotelasticsearch", elasticsearch.NewElasticsearchExecutor)
	tsdb.RegisterTsdbQueryEndpoint("gmetricsd", NewGMetricsdExecutor)
}
