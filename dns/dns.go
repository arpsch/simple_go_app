package dns

import (
	"context"
	"errors"

	"github.com/arpsch/ha/model"
)

var (
	ErrInvalidParams = errors.New("invalid parameters supplied")
)

// DNSApp represents the behavour on DNS
type DNSApp interface {
	GetLocation(ctx context.Context, pos model.Position) (float64, error)
}

// dnsApp is an the DNS object that serves the behaviors
type dnsApp struct {
	//  to add applicaiton specific properties
}

// NewApp initialize a new DNS App
func NewApp() *dnsApp {

	return &dnsApp{}
}

// GetLocation returns the location calculated based on the formula
func (dns *dnsApp) GetLocation(ctx context.Context, pos model.Position) (float64, error) {

	/* formula: loc = x*SectorID + y*SectorID + z*SectorID + vel */

	if pos.X == 0 || pos.Y == 0 || pos.Z == 0 || pos.SectorID == 0 {
		return 0.0, ErrInvalidParams
	}

	return (pos.X+pos.Y+pos.Z)*float64(pos.SectorID) + pos.Velocity, nil

}
