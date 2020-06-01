package elasticsearch

import (
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

func TestElasticGetPreferredNamesForQuery(t *testing.T) {
	Convey("Test Elasticsearch getPreferredNamesForQueries", t, func() {
		Convey("Get Name with No Alias", func() {
			testModelJSON := `
		{
					"metrics": [
						{
							"field": "value",
							"id": "1",
							"type": "avg"
						},
						{
							"field": "1",
							"id": "3",
							"pipelineAgg": "1",
							"type": "moving_avg"
						}
					]
		}`
			queries := &tsdb.Query{}
			var err error
			queries.Model, err = simplejson.NewJson([]byte(testModelJSON))
			So(err, ShouldBeNil)

			names := getPreferredNamesForQueries(queries)
			So(len(names), ShouldEqual, 2)
			So(names.GetName("3"), ShouldEqual, "Moving Average Average value")
			So(names.GetName("1"), ShouldEqual, "Average value")
			So(names.GetName("???"), ShouldEqual, "???")

		})

		Convey("Get Name with Alias", func() {
			testModelJSON := `
		{
		      "metrics": [
		        {
		          "field": "value",
		          "id": "1",
		          "type": "avg"
		        },
		        {
		          "field": "1",
		          "id": "3",
		          "pipelineAgg": "1",
		          "type": "moving_avg"
		        }
		      ],
					"alias": "overridden by alias"
		}`
			queries := &tsdb.Query{}
			var err error
			queries.Model, err = simplejson.NewJson([]byte(testModelJSON))
			So(err, ShouldBeNil)

			names := getPreferredNamesForQueries(queries)
			So(len(names), ShouldEqual, 2)
			So(names.GetName("3"), ShouldEqual, "overridden by alias")
			So(names.GetName("1"), ShouldEqual, "overridden by alias")
			So(names.GetName("???"), ShouldEqual, "???")
		})
	})
}

func TestElasticsearchGetIndexList(t *testing.T) {
	Convey("Test Elasticsearch getIndex ", t, func() {
		// TODO we cannot give pre defined now value to TimeRange anymore. It is only assigned by the constructor.

		timeRange := &tsdb.TimeRange{
			From: "48h",
			To:   "now",
		}

		Convey("Parse Interval Formats", func() {
			So(getIndex("[logstash-]YYYY.MM.DD", "Daily", timeRange),
				ShouldEqual, "logstash-2017.02.16,logstash-2017.02.17,logstash-2017.02.18")

			timeRange.From = "3h"
			So(getIndex("[logstash-]YYYY.MM.DD.HH", "Hourly", timeRange),
				ShouldEqual, "logstash-2017.02.18.09,logstash-2017.02.18.10,logstash-2017.02.18.11,logstash-2017.02.18.12")

			timeRange.From = "200h"
			So(getIndex("[logstash-]YYYY.W", "Weekly", timeRange),
				ShouldEqual, "logstash-2017.6,logstash-2017.7")

			timeRange.From = "700h"
			So(getIndex("[logstash-]YYYY.MM", "Monthly", timeRange),
				ShouldEqual, "logstash-2017.01,logstash-2017.02")

			timeRange.From = "10000h"
			So(getIndex("[logstash-]YYYY", "Yearly", timeRange),
				ShouldEqual, "logstash-2015,logstash-2016,logstash-2017")
		})

		Convey("No Interval", func() {
			index := getIndex("logstash-test", "", timeRange)
			So(index, ShouldEqual, "logstash-test")
		})
	})
}
