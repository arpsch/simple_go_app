//go:build unit
// +build unit

package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arpsch/ha/dns"
	mdns "github.com/arpsch/ha/dns/mocks"
	"github.com/arpsch/ha/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// makeMockApiHandler
func makeMockApiHandler(t *testing.T, dnsApp dns.DNSApp) http.Handler {
	handlers, err := NewRouter(dnsApp)
	assert.NotNil(t, handlers)
	assert.NoError(t, err)

	return handlers
}

func TestGetLocationHandler(t *testing.T) {

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
		err        error
		statusCode int
	}{
		{
			name:   "with valid parameters",
			method: http.MethodGet,
			URL:    "/api/v1/dns/drones/%s/location",

			want: 4929.4,

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

			want: 112,

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

			qpSecID: "4553",
			did:     "122",

			err:        dns.ErrInvalidParams,
			statusCode: http.StatusForbidden,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {

			url := fmt.Sprintf(tc.URL, tc.did)

			mockDnsApp := &mdns.MockDnsApp{}

			mockDnsApp.On("GetLocation",
				mock.AnythingOfType("*context.valueCtx"),
				mock.AnythingOfType("model.Position")).Return(tc.want, tc.err)

			// dont set as the test is against its invalid value
			if tc.name != "with invalid parameter: SectorID" {
				os.Setenv(SECTOR_ID, tc.qpSecID)
			}
			apiHandler := makeMockApiHandler(t, mockDnsApp)

			req, err := http.NewRequest(tc.method, url, nil)
			assert.Nil(t, err)

			values := req.URL.Query()
			values.Add("x", tc.qpX)
			values.Add("y", tc.qpY)
			values.Add("z", tc.qpZ)
			values.Add("velocity", tc.qpVel)
			values.Add("sector_id", tc.qpSecID)
			req.URL.RawQuery = values.Encode()

			recorder := httptest.NewRecorder()
			apiHandler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.statusCode, recorder.Result().StatusCode)
		})
	}
}

func TestParseFloat64(t *testing.T) {

	tt := []struct {
		name string

		input string
		want  float64
		err   error
	}{
		{
			name:  "with valid parameters",
			input: "4929.4",
			want:  4929.4,

			err: nil,
		},
		{
			name: "with invalid valid parameter",

			input: "1q1",
			want:  -1,

			err: dns.ErrInvalidParams,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {

			fl, err := parseFloat64(tc.input)

			assert.Equal(t, tc.want, fl)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestParseQueryParams(t *testing.T) {

	tt := []struct {
		name string

		qpX   string
		qpY   string
		qpZ   string
		qpVel string

		qpSecID string

		want model.Position
		err  error
	}{
		{
			name:  "with valid parameters",
			qpX:   "12.0",
			qpY:   "13.9",
			qpZ:   "13.9",
			qpVel: "34.0",

			qpSecID: "123",

			want: model.Position{
				X:        12.0,
				Y:        13.9,
				Z:        13.9,
				Velocity: 34.0,
				SectorID: 123,
			},

			err: nil,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", "http:/localhost:8080/api/v1?x="+tc.qpX+"&y="+tc.qpY+"&z="+tc.qpZ+"&velocity="+tc.qpVel+"&sector_id="+tc.qpSecID, nil)
			if err != nil {
				t.Fatal(err)
			}

			pos, err := parseQueryParams(req, queryParamX, queryParamY, queryParamZ, queryParamVelocity, queryParamSectorID)

			assert.Equal(t, tc.want, pos)
			assert.Equal(t, tc.err, err)
		})
	}
}
