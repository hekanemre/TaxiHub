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

// GetDriver godoc
// @Summary      Get driver by ID
// @Description  Retrieves a driver's details by their unique ID.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        id   query     string  true  "Driver ID"
// @Success      200  {object}  GetDriverResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router       /drivers/getbyid/ [get]
func (h *GetDriverHandler) Handle(ctx context.Context, req *GetDriverRequest) (*GetDriverResponse, error) {
	driver, err := h.repo.GetDriverByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetDriverResponse{
		Driver: driver,
	}, nil
}
