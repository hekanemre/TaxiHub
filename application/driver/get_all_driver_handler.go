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

// GetAllDriver godoc
// @Summary      Get all drivers
// @Description  Retrieves a paginated list of all drivers.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"       default(1)
// @Param        page_size  query     int     false  "Number of items per page" default(10)
// @Success      200  {object}  GetAllDriverResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router       /drivers/getall [get]
func (h *GetAllDriverHandler) Handle(ctx context.Context, req *GetAllFilterRequest) (*GetAllDriverResponse, error) {
	drivers, err := h.repo.GetAllDrivers(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &GetAllDriverResponse{
		Driver: drivers,
	}, nil
}
