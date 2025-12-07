package application

import (
	"context"

	"github.com/hekanemre/taxihub/domain"
)

type GetDriverByPlateHandler struct {
	repo Repository
}

type GetDriverByPlateRequest struct {
	Plate string `json:"plate"`
}

type GetDriverByPlateResponse struct {
	Driver *domain.Driver `json:"driver"`
}

func NewGetDriverByPlateHandler(repo Repository) *GetDriverByPlateHandler {
	return &GetDriverByPlateHandler{
		repo: repo,
	}
}

func (h *GetDriverByPlateHandler) Handle(ctx context.Context, req *GetDriverByPlateRequest) (*GetDriverByPlateResponse, error) {
	driver, err := h.repo.GetDriverByPlate(ctx, req.Plate)
	if err != nil {
		return nil, err
	}

	return &GetDriverByPlateResponse{
		Driver: driver,
	}, nil
}
