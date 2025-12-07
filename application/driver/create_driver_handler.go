package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hekanemre/taxihub/domain"
)

type CreateDriverHandler struct {
	repo Repository
}

type CreateDriverRequest struct {
	FirstName string          `bson:"firstName" json:"firstName"`
	LastName  string          `bson:"lastName" json:"lastName"`
	Plate     string          `bson:"plate" json:"plate"`
	TaxiType  string          `bson:"taxiType" json:"taksiType"`
	CarBrand  string          `bson:"carBrand" json:"carBrand"`
	CarModel  string          `bson:"carModel" json:"carModel"`
	Location  domain.Location `bson:"location" json:"location"`
	CreatedAt time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time       `bson:"updatedAt" json:"updatedAt"`
}

type CreateDriverResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCreateDriverHandler(repo Repository) *CreateDriverHandler {
	return &CreateDriverHandler{
		repo: repo,
	}
}

// CreateDriver godoc
// @Summary      Create a new driver
// @Description  Creates a new driver with the provided details.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        driver  body      CreateDriverRequest  true  "Driver creation data"
// @Success      200  {object}  CreateDriverResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router       /drivers/create [post]
func (h *CreateDriverHandler) Handle(ctx context.Context, req *CreateDriverRequest) (*CreateDriverResponse, error) {

	driver := &domain.Driver{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location:  req.Location,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}

	// ensure ID and timestamps are set before insert
	driver.ID = uuid.New().String()
	if driver.CreatedAt.IsZero() {
		driver.CreatedAt = time.Now()
	}
	if driver.UpdatedAt.IsZero() {
		driver.UpdatedAt = driver.CreatedAt
	}

	err := h.repo.CreateDriver(ctx, driver)
	if err != nil {
		return nil, err
	}

	return &CreateDriverResponse{
		ID:   driver.ID,
		Name: driver.FirstName + " " + driver.LastName,
	}, nil
}
