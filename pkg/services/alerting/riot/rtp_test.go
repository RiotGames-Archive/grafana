package riot

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	EXISTING_ENABLED_URL      = "http://elasticsearch5-metrics.rdatasrv.net:9200"
	EXISTING_ENABLED_URL2     = "http://elasticsearch5-metrics2.rdatasrv.net:9200/"
	EXISTING_DEGRADED_URL     = "http://elasticsearch5-metrics-w1.rdatasrv.net:9200"
	EXISTING_EMPTY_STATUS_URL = "http://elasticsearch5-metrics-w2.rdatasrv.net:9200"
	NON_EXISTING_URL          = "http://elasticsearch5-something.rdatasrv.net:9200"
	EMPTY_URL                 = ""
	NO_PORT_URL               = "http://elasticsearch5-metrics3.rdatasrv.net"
	MALFORMED_URL             = "dfsdfs"
)

func TestIsRTPHealthy(t *testing.T) {
	mockRtpHealthCheckerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
  "elasticsearch5-metrics-w1.rdatasrv.net": {
    "status": "DEGRADED",
    "reason": "",
    "details": null
  },
  "elasticsearch5-metrics-w2.rdatasrv.net": {
    "status": "",
    "reason": "",
    "details": null
  },
  "elasticsearch5-metrics.rdatasrv.net": {
    "status": "ENABLED",
    "reason": "",
    "details": null
  },
  "elasticsearch5-metrics2.rdatasrv.net": {
    "status": "ENABLED",
    "reason": "",
    "details": null
  }
}`
		w.Write([]byte(response))
	}))

	os.Setenv(RTP_CIRCUIT_BREAKER_ENDPOINT, mockRtpHealthCheckerServer.URL)
	defer os.Unsetenv(RTP_CIRCUIT_BREAKER_ENDPOINT)

	tests := []struct {
		name    string
		dsurl   string
		want    bool
		wantErr bool
	}{
		{"EMPTY_URL", EMPTY_URL, true, false},
		{"EXISTING_ENABLED_URL", EXISTING_ENABLED_URL, true, false},
		{"EXISTING_ENABLED_URL2", EXISTING_ENABLED_URL2, true, false},
		{"EXISTING_EMPTY_STATUS_URL", EXISTING_EMPTY_STATUS_URL, true, false},
		{"NON_EXISTING_URL", NON_EXISTING_URL, true, false},
		{"NO_PORT_URL", NO_PORT_URL, true, false},
		{"MALFORMED_URL", MALFORMED_URL, true, false},
		{"EXISTING_DEGRADED_URL", EXISTING_DEGRADED_URL, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health, err := IsRTPHealthy(tt.dsurl)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, health)
		})
	}

}
