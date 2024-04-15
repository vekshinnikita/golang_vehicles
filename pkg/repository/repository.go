package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/internal/filter"
	"github.com/vekshinnikita/golang_vehicles/internal/pagination"
)

type Vehicle interface {
	CreateVehicles(user []golang_vehicles.Vehicle) ([]int, error)
	UpdateVehicle(vehicleId int, vehicle golang_vehicles.UpdateVehicle) error
	GetAllVehicles(filter *filter.Filters, pagination *pagination.Pagination) (pagination.PaginatedResponse[[]golang_vehicles.Vehicle], error)
	DeleteVehicle(vehicleId int) error
}

type Repository struct {
	Vehicle
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Vehicle: NewVehiclePostgres(db),
	}
}
