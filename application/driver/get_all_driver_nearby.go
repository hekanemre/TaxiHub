package application

import (
	"context"
	"sort"
)

type GetAllDriverNearbyHandler struct {
	repo Repository
}

type GetAllDriverNearbyRequest struct {
	Lat      float64 `bson:"lat" json:"lat"`
	Lon      float64 `bson:"lon" json:"lon"`
	TaxiType string  `bson:"taxiType" json:"taxiType"`
}

type GetAllDriverNearbyResponse struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Plate      string  `json:"plate"`
	DistanceKm float64 `json:"distanceKm"`
}

func NewGetAllDriverNearbyHandler(repo Repository) *GetAllDriverNearbyHandler {
	return &GetAllDriverNearbyHandler{
		repo: repo,
	}
}

// GetAllDriverNearby godoc
// @Summary      Get all nearby drivers
// @Description  Retrieves a list of all nearby drivers based on location and taxi type.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        lat       query     float64     true  "Latitude"
// @Param        lon       query     float64     true  "Longitude"
// @Param        taxiType  query     string      true  "Type of taxi"
// @Success      200  {array}  GetAllDriverNearbyResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router       /drivers/getallnearby [get]
func (h *GetAllDriverNearbyHandler) Handle(ctx context.Context, req *GetAllDriverNearbyRequest) ([]*GetAllDriverNearbyResponse, error) {
	drivers, err := h.repo.GetAllDriversNearby(ctx, req.Lat, req.Lon, req.TaxiType)
	if err != nil {
		return nil, err
	}

	var responses []*GetAllDriverNearbyResponse
	for _, driver := range drivers {
		// MongoDB $geoNear returns distance in meters; convert to km
		distanceKm := HaversineKm(req.Lat, req.Lon, driver.Location.Coordinates[1], driver.Location.Coordinates[0])
		responses = append(responses, &GetAllDriverNearbyResponse{
			FirstName:  driver.FirstName,
			LastName:   driver.LastName,
			Plate:      driver.Plate,
			DistanceKm: distanceKm,
		})
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].DistanceKm < responses[j].DistanceKm
	})

	return responses, nil
}
