package service

import (
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/internal/filter"
	"github.com/vekshinnikita/golang_vehicles/internal/pagination"
	"github.com/vekshinnikita/golang_vehicles/pkg/repository"
)

type Vehicle interface {
	CreateVehicles(regNums []string) ([]int, error)
	UpdateVehicles(vehicleId int, vehicle golang_vehicles.UpdateVehicle) error
	GetAllVehicles(filter *filter.Filters, pagination *pagination.Pagination) (pagination.PaginatedResponse[[]golang_vehicles.Vehicle], error)
	DeleteVehicle(vehicleId int) error
}

type Service struct {
	Vehicle
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Vehicle: NewVehicleService(repo.Vehicle),
	}
}
