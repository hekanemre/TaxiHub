package application

import (
	"context"

	"github.com/hekanemre/taxihub/domain"
)

// we add this to lose coupling. We used dependency inversion.
// so app is not directly dependent to repository
type Repository interface {
	CreateDriver(ctx context.Context, driver *domain.Driver) error
	UpdateDriver(ctx context.Context, driver *domain.Driver) error
	GetAllDrivers(ctx context.Context, page, pageSiz int) ([]*domain.Driver, error)
	GetDriverByID(ctx context.Context, id string) (*domain.Driver, error)
	GetDriverByPlate(ctx context.Context, plate string) (*domain.Driver, error)
	GetAllDriversNearby(ctx context.Context, lat, lon float64, taxiType string) ([]*domain.Driver, error)
}
