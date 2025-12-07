package application

import (
	"context"
	"time"

	"github.com/hekanemre/taxihub/domain"
)

type UpdateDriverHandler struct {
	repo Repository
}

type UpdateDriverRequest struct {
	ID        string          `bson:"id" json:"id"`
	FirstName string          `bson:"firstName" json:"firstName"`
	LastName  string          `bson:"lastName" json:"lastName"`
	Plate     string          `bson:"plate" json:"plate"`
	TaxiType  string          `bson:"taxiType" json:"taksiType"`
	CarBrand  string          `bson:"carBrand" json:"carBrand"`
	CarModel  string          `bson:"carModel" json:"carModel"`
	Location  domain.Location `bson:"location" json:"location"`
}

type UpdateDriverResponse struct {
	Driver *domain.Driver `json:"driver"`
}

func NewUpdateDriverHandler(repo Repository) *UpdateDriverHandler {
	return &UpdateDriverHandler{
		repo: repo,
	}
}

func (h *UpdateDriverHandler) Handle(ctx context.Context, req *UpdateDriverRequest) (*UpdateDriverResponse, error) {
	driver := &domain.Driver{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location:  req.Location,
	}

	driver.UpdatedAt = time.Now()
	err := h.repo.UpdateDriver(ctx, driver)
	if err != nil {
		return nil, err
	}

	return &UpdateDriverResponse{
		Driver: driver,
	}, nil
}
