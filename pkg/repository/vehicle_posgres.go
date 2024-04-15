package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/internal/filter"
	. "github.com/vekshinnikita/golang_vehicles/internal/pagination"
	"github.com/vekshinnikita/golang_vehicles/internal/tools"
)

type VehiclePostgres struct {
	db *sqlx.DB
}

func NewVehiclePostgres(db *sqlx.DB) *VehiclePostgres {
	return &VehiclePostgres{db}
}

func (r *VehiclePostgres) createOrGetUser(user golang_vehicles.Owner) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1 AND surname=$2 AND patronymic=$3", usersTable)
	err := r.db.Get(&id, query, user.Name, user.Surname, user.Patronymic)

	if err == nil {
		return id, nil
	}

	query = fmt.Sprintf("INSERT INTO %s (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Patronymic)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *VehiclePostgres) CreateVehicles(vehicles []golang_vehicles.Vehicle) ([]int, error) {
	userIds := make([]int, 0)

	for _, vehicle := range vehicles {
		userId, err := r.createOrGetUser(vehicle.Owner)
		if err != nil {
			return make([]int, 0), err
		}
		userIds = append(userIds, userId)
	}

	argsCounter := 1
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	for index, vehicle := range vehicles {
		setValues = append(setValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", argsCounter, argsCounter+1, argsCounter+2, argsCounter+3, argsCounter+4))
		args = append(args, vehicle.Model, vehicle.Mark, vehicle.RegNum, vehicle.Year, userIds[index])
		argsCounter += 5
	}

	setQueries := strings.Join(setValues, ", ")

	var ids []int
	query := fmt.Sprintf("INSERT INTO %s (model, mark, reg_num, year, owner_id) VALUES %s RETURNING id", vehiclesTable, setQueries)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return make([]int, 0), err
	}

	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		ids = append(ids, id)
	}

	//Info log
	stringArgs, _ := tools.Map(args, func(v any) (any, error) { return fmt.Sprint(v), nil })

	re := regexp.MustCompile(`\$\d`)
	infoQueryString := string(re.ReplaceAllString(setQueries, "%s"))
	infoQueryString = fmt.Sprintf(infoQueryString, stringArgs...)

	logrus.Info(fmt.Sprintf("record have been created (model, mark, reg_num, year, owner_id) VALUES %s", infoQueryString))

	return ids, nil
}

func (r *VehiclePostgres) getVehicleById(vehicleId int) (golang_vehicles.Vehicle, error) {
	var vehicle golang_vehicles.Vehicle
	query := fmt.Sprintf(`
		SELECT v.id, v.mark, v.model, v.reg_num, v.year, u.id, u.name, u.surname, u.patronymic 
			FROM %s as v 
			JOIN %s as u 
				ON v.owner_id=u.id 
			WHERE v.id=$1
	`, vehiclesTable, usersTable)

	row := r.db.QueryRow(query, vehicleId)
	err := row.Scan(
		&vehicle.VehicleId,
		&vehicle.Mark,
		&vehicle.Model,
		&vehicle.RegNum,
		&vehicle.Year,
		&vehicle.Owner.OwnerId,
		&vehicle.Owner.Name,
		&vehicle.Owner.Surname,
		&vehicle.Owner.Patronymic,
	)

	if err == nil {
		return vehicle, nil
	}

	return vehicle, nil
}

func (r *VehiclePostgres) UpdateVehicle(vehicleId int, updateVehicle golang_vehicles.UpdateVehicle) error {

	vehicle, _ := r.getVehicleById(vehicleId)

	setValues, args, err := tools.GetUpdateQueryString(updateVehicle, []string{"VehicleId", "Owner"})
	if err != nil {
		return err
	}

	if (updateVehicle.Owner != golang_vehicles.UpdateOwner{}) {
		owner := vehicle.Owner

		if updateVehicle.Owner.Name != "" {
			owner.Name = updateVehicle.Owner.Name
		}
		if updateVehicle.Owner.Surname != "" {
			owner.Surname = updateVehicle.Owner.Surname
		}
		if updateVehicle.Owner.Patronymic != "" {
			owner.Patronymic = updateVehicle.Owner.Patronymic
		}

		owner_id, err := r.createOrGetUser(owner)
		if err != nil {
			return err
		}

		args = append(args, owner_id)
		setValues = append(setValues, fmt.Sprintf("owner_id=$%d", len(args)))
	}

	setQueries := strings.Join(setValues, ", ")
	args = append(args, vehicleId)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", vehiclesTable, setQueries, len(args))
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	//Info log
	stringArgs, _ := tools.Map(args[:len(args)-1], func(v any) (any, error) { return fmt.Sprint(v), nil })

	re := regexp.MustCompile(`\$\d`)
	infoQueryString := string(re.ReplaceAllString(setQueries, "%s"))
	infoQueryString = fmt.Sprintf(infoQueryString, stringArgs...)

	logrus.Info(fmt.Sprintf("the fields of the record with id=%d have been updated to %s", vehicleId, infoQueryString))

	return nil
}

func (r *VehiclePostgres) GetAllVehicles(filter *filter.Filters, pagination *Pagination) (PaginatedResponse[[]golang_vehicles.Vehicle], error) {
	var vehicles []golang_vehicles.Vehicle

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), v.id as id, v.mark, v.model, v.reg_num as regNum, v.year, owner.id as owner_id, owner.name as owner_name, owner.surname as owner_surname, owner.patronymic as owner_patronymic
		FROM %s as v
		JOIN %s as owner
			ON v.owner_id=owner.id 
		%s
		%s
		%s
	`, vehiclesTable, usersTable, filter.GetSQLFilterQuery(), pagination.GetSQLSortQuery(), pagination.GetSQLPaginationQuery())

	fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return PaginatedResponse[[]golang_vehicles.Vehicle]{}, err
	}

	var totalRecords int

	for rows.Next() {
		var vehicle golang_vehicles.Vehicle
		rows.Scan(
			&totalRecords,
			&vehicle.VehicleId,
			&vehicle.Mark,
			&vehicle.Model,
			&vehicle.RegNum,
			&vehicle.Year,
			&vehicle.Owner.OwnerId,
			&vehicle.Owner.Name,
			&vehicle.Owner.Surname,
			&vehicle.Owner.Patronymic,
		)

		vehicles = append(vehicles, vehicle)
	}

	return GetPaginatedResponse(*pagination, vehicles, totalRecords), nil
}

func (r *VehiclePostgres) DeleteVehicle(vehicleId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, vehiclesTable)

	_, err := r.db.Exec(query, vehicleId)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("vehicle record with id=%d is deleted", vehicleId))

	return nil
}
