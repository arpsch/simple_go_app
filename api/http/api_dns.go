package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/arpsch/ha/dns"
	"github.com/arpsch/ha/middleware"
	"github.com/arpsch/ha/model"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

const (
	// query param constants
	queryParamX        = "x"
	queryParamY        = "y"
	queryParamZ        = "z"
	queryParamVelocity = "velocity"
	queryParamSectorID = "sector_id"

	// path param constants
	pathParamDroneID = "did"

	// SECTOR_ID holds the environment value set
	SECTOR_ID = "SECTOR_ID"
	float_64  = 64
)

type ApiHandler struct {
	App dns.DNSApp
}

func NewApiHandler(app dns.DNSApp) *ApiHandler {
	return &ApiHandler{
		App: app,
	}
}

func NewRouter(app dns.DNSApp) (*mux.Router, error) {
	apiHandler := NewApiHandler(app)

	router := mux.NewRouter()

	// Device API
	router.HandleFunc("/api/v1/dns/drones/{did:[0-9]+}/location", apiHandler.GetLocationHandler).Methods("GET")

	// internal API
	router.HandleFunc("/api/v1/internal/dns/drones/{did:[0-9]+}/location", apiHandler.GetLocationHandler).Methods("GET")

	router.Use(middleware.Logger)

	// when sector id is not set, could be used for multi sectors
	sectorID := os.Getenv(SECTOR_ID)
	if sectorID != "" {
		secValidator, err := middleware.NewSectorIDValidator(queryParamSectorID, sectorID)
		if err != nil {
			return nil, err
		}

		router.Use(secValidator.Middleware)
	}

	return router, nil
}

func parseFloat64(str string) (float64, error) {
	if str == "" {
		return -1, errors.New("emtpy string")
	}

	fl, err := strconv.ParseFloat(str, float_64)
	if err != nil {
		return -1, dns.ErrInvalidParams
	}
	return fl, nil
}

func parseQueryParams(r *http.Request, queryParams ...string) (model.Position, error) {
	pos := model.Position{}

	qpVals := make(map[string]string)

	for _, qp := range queryParams {
		qpVals[qp] = r.URL.Query().Get(qp)
		fl, err := parseFloat64(qpVals[qp])
		if err != nil {
			return model.Position{}, err
		}

		switch qp {
		case queryParamX:
			pos.X = fl
		case queryParamY:
			pos.Y = fl
		case queryParamZ:
			pos.Z = fl
		case queryParamVelocity:
			pos.Velocity = fl
		case queryParamSectorID:
			pos.SectorID = int(fl)
		}
	}

	return pos, nil
}

// GetLocationHandler handles the get location request from the drones
// returns location based on the parameters
func (ah *ApiHandler) GetLocationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	did := vars["did"]
	_, err := strconv.Atoi(vars["did"])
	if err != nil {
		http.Error(w,
			errors.Wrap(err, "failed to parse the payload").Error(),
			http.StatusBadRequest)
		return
	}

	log.Printf("request from drone %s recieved", did)

	pos, err := parseQueryParams(r,
		queryParamX,
		queryParamY,
		queryParamZ,
		queryParamVelocity,
		queryParamSectorID,
	)
	if err != nil {
		http.Error(w,
			errors.Wrap(err, "failed to parse the payload").Error(),
			http.StatusBadRequest)
		return
	}

	// if needed for future or logging
	pos.DroneID = did

	loc, err := ah.App.GetLocation(ctx, pos)
	if err != nil {
		if errors.Is(err, dns.ErrInvalidParams) {
			http.Error(w,
				errors.Wrapf(err, "failed to retrieve location: %v", did).Error(),
				http.StatusBadRequest)
			return

		}

		http.Error(w,
			"failed retrieve location:"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	locRes := Location{
		Loc: &loc,
	}

	locJsonRes, err := json.Marshal(&locRes)
	if err != nil {

		http.Error(w,
			errors.Wrapf(err, "failed to retrieve location: %v", did).Error(),
			http.StatusInternalServerError)
		return
	}

	w.Write(locJsonRes)
}
