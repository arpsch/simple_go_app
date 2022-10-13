package middleware

import (
	"errors"
	"net/http"
)

// SectorIDValidator represents sector validator
type SectorIDValidator struct {
	Key   string
	Value string
}

// NewSectorIDValidator constructs the validator object with sector id
// and its key
func NewSectorIDValidator(secIDKey, secID string) (*SectorIDValidator, error) {

	if secIDKey == "" {
		return nil, errors.New("missing param")
	}

	return &SectorIDValidator{
		Key:   secIDKey,
		Value: secID,
	}, nil
}

// Middleware function, which will be called for each request
func (sid *SectorIDValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sectorID := r.URL.Query().Get(sid.Key)

		// if the sector id is not set, this instance could be used for
		// multiple sectors, else for only the set sector
		if sid.Value == "" || sectorID == sid.Value {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "forbidden: requesting for wrong dns instance",
				http.StatusForbidden)
		}
	})
}
