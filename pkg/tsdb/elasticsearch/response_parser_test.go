package elasticsearch

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testResponseJSON = `
{
  "took": 1588,
  "timed_out": false,
  "_shards": {
    "total": 2250,
    "successful": 2250,
    "failed": 0
  },
  "hits": {
    "total": 5,
    "max_score": 0,
    "hits": []
  },
  "aggregations": {
    "2": {
      "buckets": [
        {
          "key_as_string": "1487020955000",
          "key": 1487020955000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487020985000",
          "key": 1487020985000,
          "doc_count": 1,
          "1": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021010000",
          "key": 1487021010000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021045000",
          "key": 1487021045000,
          "doc_count": 1,
          "1": {
            "value": 1234
          },
          "3": {
            "value": 123
          }
        },
        {
          "key_as_string": "1487021105000",
          "key": 1487021105000,
          "doc_count": 1,
          "1": {
            "value": 155
          },
          "3": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021165000",
          "key": 1487021165000,
          "doc_count": 1,
          "1": {
            "value": 0
          },
          "3": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021180000",
          "key": 1487021180000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021210000",
          "key": 1487021210000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021225000",
          "key": 1487021225000,
          "doc_count": 1,
          "1": {
            "value": 0
          },
          "3": {
            "value": 1000
          }
        }
      ]
    }
  }
}`

var testRecursiveResponseJSON = `
{
  "aggregations": {
    "2": {
      "buckets": [
      {
        "3": {
          "buckets": [
          { "4": {"value": 10}, "doc_count": 1, "key": 1000},
          { "4": {"value": 12}, "doc_count": 3, "key": 2000}
          ]
        },
        "doc_count": 4,
        "key": "server1"
      },
      {
        "3": {
          "buckets": [
          { "4": {"value": 20}, "doc_count": 1, "key": 1000},
          { "4": {"value": 32}, "doc_count": 3, "key": 2000}
          ]
        },
        "doc_count": 10,
        "key": "server2"
      }
      ]
    }
  }
}`

var testWildcardResponseJSON = `
{
  "took": 707,
  "timed_out": false,
  "_shards": {
    "total": 20,
    "successful": 20,
    "failed": 0
  },
  "hits": {
    "total": 14358,
    "max_score": 0.0,
    "hits": []
  },
  "aggregations": {
    "3": {
      "doc_count_error_upper_bound": 0,
      "sum_other_doc_count": 0,
      "buckets": [
        {
          "key": "lolgarena.jkt1.id1_honor1.honor.server",
          "doc_count": 900,
          "2": {
            "buckets": [
              {
                "key_as_string": "1505922780000",
                "key": 1505922780000,
                "doc_count": 5,
                "1": {
                  "value": 5.0
                }
              },
              {
                "key_as_string": "1505922840000",
                "key": 1505922840000,
                "doc_count": 5,
                "1": {
                  "value": 5.0
                }
              },
              {
                "key_as_string": "1505922842000",
                "key": 1505922840000,
                "doc_count": 5,
                "1": {
                  "value": 5.0
                }
              }
            ]
          }
        },
        {
          "key": "lolgarena.ph1.ph_honor1.honor.server",
          "doc_count": 900,
          "2": {
            "buckets": [
              {
                "key_as_string": "1505912820000",
                "key": 1505912820000,
                "doc_count": 0,
                "1": {
                  "value": 0.0
                }
              },
              {
                "key_as_string": "1505923560000",
                "key": 1505923560000,
                "doc_count": 5,
                "1": {
                  "value": 5.0
                }
              },
              {
                "key_as_string": "1505923620000",
                "key": 1505923620000,
                "doc_count": 5,
                "1": {
                  "value": 5.0
                }
              }
            ]
          }
        }
      ]
    }
  },
  "status": 200
}
`

var testResponseNullDataJSON = `
{
  "took": 1588,
  "timed_out": false,
  "_shards": {
    "total": 2250,
    "successful": 2250,
    "failed": 0
  },
  "hits": {
    "total": 0,
    "max_score": 0,
    "hits": []
  },
  "aggregations": {
    "2": {
      "buckets": [
        {
          "key_as_string": "1487020955000",
          "key": 1487020955000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021010000",
          "key": 1487021010000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021180000",
          "key": 1487021180000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021210000",
          "key": 1487021210000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        }
      ]
    }
  }
}`

var testMultiBucketResponse = `{"took":36,"timed_out":false,"_shards":{"total":20,"successful":20,"failed":0},"hits":{"total":114,"max_score":0.0,"hits":[]},"aggregations":{"4":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"builtin.general.enabled_instance_count","doc_count":114,"6":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"las2","doc_count":12,"7":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"cacti","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"cauth","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"nxlog","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}}]}},{"key":"pdx3","doc_count":8,"7":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"monitor","doc_count":8,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":2,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}}]}},{"key":"pdx32","doc_count":62,"7":{"doc_count_error_upper_bound":0,"sum_other_doc_count":7,"buckets":[{"key":"cacti","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"cauth","doc_count":7,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.5}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.6666666666666667}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.5}}]}},{"key":"dns","doc_count":6,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":1.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":0.5}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":1.0},"5":{"value":0.6666666666666666}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"exabgp","doc_count":6,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":2.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":2.0},"5":{"value":2.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":2.0},"5":{"value":2.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"maas","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"ntp","doc_count":6,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.3333333333333333}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"nxlog","doc_count":7,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.5}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":2.0},"5":{"value":1.6666666666666667}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"proxy","doc_count":7,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":2,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"pulp","doc_count":6,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":2,"1":{"value":2.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":2,"1":{"value":2.0},"5":{"value":2.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":2,"1":{"value":2.0},"5":{"value":2.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"rims","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}}]}},{"key":"pdx33","doc_count":32,"7":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"cacti","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"cauth","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"dns","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"exabgp","doc_count":4,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":0.0},"5":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"maas","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"ntp","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"nxlog","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"proxy","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"pulp","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}},{"key":"tacacs","doc_count":3,"2":{"buckets":[{"key_as_string":"1506613590000","key":1506613590000,"doc_count":0,"1":{"value":0.0}},{"key_as_string":"1506613620000","key":1506613620000,"doc_count":1,"1":{"value":1.0}},{"key_as_string":"1506613650000","key":1506613650000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613680000","key":1506613680000,"doc_count":1,"1":{"value":1.0},"5":{"value":1.0}},{"key_as_string":"1506613710000","key":1506613710000,"doc_count":0,"1":{"value":0.0}}]}}]}}]}}]}},"status":200}`

func TestElasticserachQueryParser(t *testing.T) {
	Convey("Elasticserach QueryBuilder query parsing", t, func() {

		Convey("Parse ElasticSearch Query Results", func() {
			names := NameMap{}
			names["1"] = Name{
				Value: "Average value",
			}
			names["3"] = Name{
				Value:     "Moving Average",
				Reference: "1",
			}
			queryResult, err := parseQueryResult([]byte(testResponseJSON), names, FilterMap{})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)

			So(len(queryResult.Series), ShouldEqual, 2)
		})

		Convey("Parse ElasticSearch Nested Query Results", func() {
			names := NameMap{}
			queryResult, err := parseQueryResult([]byte(testRecursiveResponseJSON), names, FilterMap{})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)
			So(len(queryResult.Series), ShouldEqual, 2)
			So(queryResult.Series[0].Name, ShouldContainSubstring, "server")
		})

		Convey("Parse ElasticSearch Nested Query Results With Filter", func() {
			names := NameMap{}
			queryResult, err := parseQueryResult([]byte(testRecursiveResponseJSON), names, FilterMap{"4": true})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)
			So(len(queryResult.Series), ShouldEqual, 2)
		})

		Convey("Parse ElasticSearch All Null Results", func() {
			names := NameMap{}
			queryResult, err := parseQueryResult([]byte(testResponseNullDataJSON), names, FilterMap{})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)

			So(len(queryResult.Series), ShouldEqual, 0)
		})
	})
}

func TestElasticserachWildcardQueryParser(t *testing.T) {
	Convey("Elasticserach QueryBuilder WildCard query parsing", t, func() {

		Convey("Parse ElasticSearch WildCard Query Results", func() {
			names := NameMap{}
			queryResult, err := parseQueryResult([]byte(testWildcardResponseJSON), names, FilterMap{})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)
			So(len(queryResult.Series), ShouldEqual, 2)
		})
	})
}

func TestElasticserachMultiBucketResponseParser(t *testing.T) {
	Convey("TestElasticserachMultiBucketResponseParser", t, func() {

		Convey("Parse ElasticSearch Multi Bucket Results", func() {
			names := NameMap{}
			queryResult, err := parseQueryResult([]byte(testMultiBucketResponse), names, FilterMap{})

			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)
			So(len(queryResult.Series), ShouldEqual, 29)
		})
	})
}
