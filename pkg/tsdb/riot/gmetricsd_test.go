package riot

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testGmetricsdResponse = `{
  "Metadata": {
   "ResponseTime": "28.83327ms"
  },
  "History": {
   "Uuid": "222495028366932659",
   "TimeSeries": {
    "1487342347": {
     "StartTime": "1487342347",
     "EndTime": "1487342347",
     "Min": {
      "Timestamp": "1487342347",
      "Value": 0.0726
     },
     "Max": {
      "Timestamp": "1487342347",
      "Value": 0.0726
     },
     "Average": 0.0726,
     "Midrange": 0.0726,
     "NumberOfDataPoints": 1
    },
    "1487342423": {
     "StartTime": "1487342423",
     "EndTime": "1487342423",
     "Min": {
      "Timestamp": "1487342423",
      "Value": 0.0815
     },
     "Max": {
      "Timestamp": "1487342423",
      "Value": 0.0815
     },
     "Average": 0.0815,
     "Midrange": 0.0815,
     "NumberOfDataPoints": 1
    },
    "1487342502": {
     "StartTime": "1487342502",
     "EndTime": "1487342502",
     "Min": {
      "Timestamp": "1487342502",
      "Value": 0.1001
     },
     "Max": {
      "Timestamp": "1487342502",
      "Value": 0.1001
     },
     "Average": 0.1001,
     "Midrange": 0.1001,
     "NumberOfDataPoints": 1
    },
    "1487342580": {
     "StartTime": "1487342580",
     "EndTime": "1487342580",
     "Min": {
      "Timestamp": "1487342580",
      "Value": 0.0751
     },
     "Max": {
      "Timestamp": "1487342580",
      "Value": 0.0751
     },
     "Average": 0.0751,
     "Midrange": 0.0751,
     "NumberOfDataPoints": 1
    },
    "1487342657": {
     "StartTime": "1487342657",
     "EndTime": "1487342657",
     "Min": {
      "Timestamp": "1487342657",
      "Value": 0.0566
     },
     "Max": {
      "Timestamp": "1487342657",
      "Value": 0.0566
     },
     "Average": 0.0566,
     "Midrange": 0.0566,
     "NumberOfDataPoints": 1
    },
    "1487342738": {
     "StartTime": "1487342738",
     "EndTime": "1487342738",
     "Min": {
      "Timestamp": "1487342738",
      "Value": 0.0635
     },
     "Max": {
      "Timestamp": "1487342738",
      "Value": 0.0635
     },
     "Average": 0.0635,
     "Midrange": 0.0635,
     "NumberOfDataPoints": 1
    },
    "1487342818": {
     "StartTime": "1487342818",
     "EndTime": "1487342818",
     "Min": {
      "Timestamp": "1487342818",
      "Value": 0.1037
     },
     "Max": {
      "Timestamp": "1487342818",
      "Value": 0.1037
     },
     "Average": 0.1037,
     "Midrange": 0.1037,
     "NumberOfDataPoints": 1
    },
    "1487342897": {
     "StartTime": "1487342897",
     "EndTime": "1487342897",
     "Min": {
      "Timestamp": "1487342897",
      "Value": 0.0557
     },
     "Max": {
      "Timestamp": "1487342897",
      "Value": 0.0557
     },
     "Average": 0.0557,
     "Midrange": 0.0557,
     "NumberOfDataPoints": 1
    }
   },
   "PrettyName": "grpavg:TELEMETRY.ALERTEROUS.pipeline.watcher.create.rate1",
   "Name": "grpavg[\"metricsd_globalriot.las2.alerterous1_TELEMETRY.ALERTEROUS\",\"metricsd[\\\"TELEMETRY.ALERTEROUS.pipeline.watcher.create.rate1\\\"]\", avg, 90s]",
   "Host": {
    "Hostname": "metricsd_globalriot.las2.alerterous1_TELEMETRY.ALERTEROUS",
    "Hostid": "11447"
   }
  }
 }`

func TestGmetricsdParseResponse(t *testing.T) {
	Convey("Test Parse GMetricsd Response", t, func() {
		Convey("Parse Response into tsdb Response", func() {
			history := &GMetricsdHistoryResponse{}

			err := json.Unmarshal([]byte(testGmetricsdResponse), history)
			So(err, ShouldBeNil)

			So(history.History.Host.Hostname, ShouldEqual, "metricsd_globalriot.las2.alerterous1_TELEMETRY.ALERTEROUS")
			So(history.History.PrettyName, ShouldEqual, "grpavg:TELEMETRY.ALERTEROUS.pipeline.watcher.create.rate1")
			So(len(history.History.TimeSeries), ShouldEqual, 8)

			tsdbQueryResult, err := history.ToTsdbQueryResult()
			So(tsdbQueryResult, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(len(tsdbQueryResult.Series[0].Points), ShouldEqual, 8)
			So(tsdbQueryResult.Series[0].Name, ShouldEqual, fmt.Sprintf("%s [%s]", history.History.PrettyName, history.History.Host.Hostname))
		})
	})
}
