//go:build integration
// +build integration

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/arpsch/ha/dns"
)

// default server and port
var (
	server = "localhost"
	port   = "8888"
)

func TestMain(m *testing.M) {

	srv := os.Getenv("SERVER")
	if srv != "" {
		server = srv
	}

	p := os.Getenv("PORT")
	if p != "" {
		port = p
	}

	os.Exit(m.Run())
}

// TestGetLocationHandler calls the test over the wire
// **Note: we could also call the direct handler APIs but
// the downside is middleware testing will be skipped as part
// of the integration test
func TestGetLocationHandlerIntegration(t *testing.T) {

	tt := []struct {
		name   string
		method string
		URL    string

		qpX   string
		qpY   string
		qpZ   string
		qpVel string

		qpSecID string
		did     string

		want       float64
		statusCode int
		err        error
	}{
		{
			name:   "with valid parameters",
			method: http.MethodGet,
			URL:    "/api/v1/dns/drones/%s/location",
			want:   4929.4,

			qpX:   "12.0",
			qpY:   "13.9",
			qpZ:   "13.9",
			qpVel: "34.0",

			qpSecID: "123",
			did:     "122",

			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name:   "with invalid parameter: X",
			method: http.MethodGet,
			URL:    "/api/v1/dns/drones/%s/location",
			want:   112,

			qpX:   "0",
			qpY:   "13.9",
			qpZ:   "13.9",
			qpVel: "34.0",

			qpSecID: "123",
			did:     "122",

			err:        dns.ErrInvalidParams,
			statusCode: http.StatusBadRequest,
		},
		{
			name:   "with invalid parameter: SectorID",
			method: http.MethodGet,
			URL:    "/api/v1/dns/drones/%s/location",
			want:   112,

			qpX:   "0",
			qpY:   "13.9",
			qpZ:   "13.9",
			qpVel: "34.0",

			qpSecID: "3342",
			did:     "122",

			statusCode: http.StatusForbidden,
		},
	}

	domain := "http://" + server + ":" + port

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			url := domain + fmt.Sprintf(tc.URL, tc.did)

			client := &http.Client{}
			req, err := http.NewRequest(tc.method, url, nil)
			if err != nil {
				t.Fatalf("failed request creation: %v\n", err)
			}

			values := req.URL.Query()
			values.Add("x", tc.qpX)
			values.Add("y", tc.qpY)
			values.Add("z", tc.qpZ)
			values.Add("velocity", tc.qpVel)
			values.Add("sector_id", tc.qpSecID)

			req.URL.RawQuery = values.Encode()

			resp, err := client.Do(req)
			if tc.err == nil && err != nil {
				t.Fatalf("expected error %v, got %v", tc.err, err)
			}

			if resp == nil {
				t.FailNow()
			}

			if tc.statusCode != resp.StatusCode {
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				t.Fatalf("expected status code %d, got %d err:%s", tc.statusCode, resp.StatusCode, bodyBytes)
			}

			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response: %v\n", err)
			}
			var loc Location
			json.Unmarshal(bodyBytes, &loc)
			if loc.Loc != nil && *loc.Loc != tc.want {
				t.Fatalf("expected location %v, got %v", tc.want, loc.Loc)
			}
		})
	}
}
