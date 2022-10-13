//go:build unit
// +build unit

package dns

import (
	"context"
	"os"
	"testing"

	"github.com/arpsch/ha/model"
)

var droneNSApp *dnsApp

func TestMain(m *testing.M) {

	droneNSApp = NewApp()

	os.Exit(m.Run())
}

func TestGetLocation(t *testing.T) {
	cases := []struct {
		name   string
		input  model.Position
		output float64
		err    error
	}{
		{
			name: "with valid input parameters",
			input: model.Position{
				X:        123.0,
				Y:        342.1,
				Z:        451.1,
				Velocity: 20.1,
				SectorID: 123,
			},
			output: 112712.70000000001,
			err:    nil,
		},
		{
			name: "with invalid input parameters: Z",
			input: model.Position{
				X:        123.0,
				Y:        342.1,
				Z:        0,
				Velocity: 20.1,
				SectorID: 123,
			},
			output: 0,
			err:    ErrInvalidParams,
		},
		{
			name: "with invalid input parameters: SectorID",
			input: model.Position{
				X:        123.0,
				Y:        342.1,
				Z:        1213,
				Velocity: 20.1,
				SectorID: 0,
			},
			output: 0,
			err:    ErrInvalidParams,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			loc, err := droneNSApp.GetLocation(context.Background(), tc.input)

			if tc.err == nil && err != nil {
				t.Fatalf("test failed with err: %v", err)
			}

			if tc.err == nil && tc.output != loc {
				t.Fatalf("expected location %v, got %v", tc.output, loc)
			}

			if tc.err != err {
				t.Fatalf("expected location %v, got %v", tc.err, err)
			}
		})
	}
}
