package model

import validation "github.com/go-ozzo/ozzo-validation"

// Position represents the current drone position with
// with velocity
type Position struct {
	DroneID string `json:"drone_id"`

	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`

	Velocity float64 `json:"velocity"`
	SectorID int     `json:"sector_id"`
}

// Validate validates the payload from the request
func (pos Position) Validate() error {
	err := validation.ValidateStruct(&pos,
		validation.Field(&pos.X, validation.Required),
		validation.Field(&pos.Y, validation.Required),
		validation.Field(&pos.Z, validation.Required),
		validation.Field(&pos.Velocity, validation.Required),
		validation.Field(&pos.SectorID, validation.Required),
	)
	if err != nil {
		return err
	}
	return nil
}
