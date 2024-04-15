package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/internal/filter"
	"github.com/vekshinnikita/golang_vehicles/internal/pagination"
	"github.com/vekshinnikita/golang_vehicles/pkg/repository"
	"golang.org/x/sync/errgroup"
)

type VehicleService struct {
	repo repository.Vehicle
}

func NewVehicleService(repo repository.Vehicle) *VehicleService {
	return &VehicleService{
		repo,
	}
}

func getVehicleInfo(regNum string, c chan golang_vehicles.Vehicle) error {
	url := fmt.Sprintf("%s?regNum=%s", viper.GetString("internal.vehicleUrl"), regNum)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var vehicle golang_vehicles.Vehicle
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&vehicle); err != nil {
		return err
	}
	c <- vehicle

	return nil
}

func (s *VehicleService) CreateVehicles(regNums []string) ([]int, error) {
	var vehicles []golang_vehicles.Vehicle
	c := make(chan golang_vehicles.Vehicle, len(regNums))
	errGrp, _ := errgroup.WithContext(context.Background())

	for _, regNum := range regNums {
		func(val string) {
			errGrp.Go(func() error { return getVehicleInfo(val, c) })
		}(regNum)
	}

	err := errGrp.Wait()
	if err != nil {
		return make([]int, 0), err
	}

	for range regNums {
		vehicles = append(vehicles, <-c)
	}

	return s.repo.CreateVehicles(vehicles)
}

func (s *VehicleService) UpdateVehicles(vehicleId int, vehicle golang_vehicles.UpdateVehicle) error {
	return s.repo.UpdateVehicle(vehicleId, vehicle)
}

func (s *VehicleService) GetAllVehicles(filter *filter.Filters, pagination *pagination.Pagination) (pagination.PaginatedResponse[[]golang_vehicles.Vehicle], error) {
	return s.repo.GetAllVehicles(filter, pagination)
}

func (s *VehicleService) DeleteVehicle(vehicleId int) error {
	return s.repo.DeleteVehicle(vehicleId)
}
