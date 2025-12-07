package application

import (
	"context"

	"github.com/hekanemre/taxihub/domain"
)

type GetAllDriverHandler struct {
	repo Repository
}

type GetAllFilterRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type GetAllDriverResponse struct {
	Driver []*domain.Driver `json:"drivers"`
}

func NewGetAllDriverHandler(repo Repository) *GetAllDriverHandler {
	return &GetAllDriverHandler{
		repo: repo,
	}
}

func (h *GetAllDriverHandler) Handle(ctx context.Context, req *GetAllFilterRequest) (*GetAllDriverResponse, error) {
	drivers, err := h.repo.GetAllDrivers(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &GetAllDriverResponse{
		Driver: drivers,
	}, nil
}
