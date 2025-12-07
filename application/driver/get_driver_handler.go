package application

import (
	"context"

	"github.com/hekanemre/taxihub/domain"
)

type GetDriverHandler struct {
	repo Repository
}

type GetDriverRequest struct {
	ID string `json:"id"`
}

type GetDriverResponse struct {
	Driver *domain.Driver `json:"driver"`
}

func NewGetDriverHandler(repo Repository) *GetDriverHandler {
	return &GetDriverHandler{
		repo: repo,
	}
}

func (h *GetDriverHandler) Handle(ctx context.Context, req *GetDriverRequest) (*GetDriverResponse, error) {
	driver, err := h.repo.GetDriverByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetDriverResponse{
		Driver: driver,
	}, nil
}
