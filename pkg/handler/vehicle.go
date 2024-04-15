package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/internal/filter"
	"github.com/vekshinnikita/golang_vehicles/internal/pagination"
	"github.com/vekshinnikita/golang_vehicles/internal/validators"
)

type createVehicleInput struct {
	RegNums []string `json:"regNums" binding:"required"`
}

// @Summary      Create vehicles
// @Tags         vehicle
// @Accept       json
// @Produce      json
// @Param        input    body     handler.createVehicleInput  true  "regNums"
// @Success      201  {array}   integer
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /vehicles [post]
func (h *Handler) CreateVehicle(c *gin.Context) {
	var input createVehicleInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(input.RegNums) <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "regNums must have at least one element")
		return
	}

	validator := validators.NewValidator()

	for index, regNum := range input.RegNums {
		validators.ValidateRegNum(validator)(fmt.Sprintf("regNums[%d]", index), regNum)
	}

	if !validator.IsValid() {
		c.JSON(http.StatusBadRequest, validator.GetErrorsDict())
		return
	}

	id, err := h.services.Vehicle.CreateVehicles(input.RegNums)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Update vehicles
// @Tags         vehicle
// @Accept       json
// @Produce      json
// @Param        input    body    golang_vehicles.UpdateVehicle  true  "Partial vehicle fields"
// @Success      204
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /vehicles/{id} [patch]
func (h *Handler) UpdateVehicle(c *gin.Context) {
	var input golang_vehicles.UpdateVehicle

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	vehicleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Vehicle.UpdateVehicles(vehicleId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Create vehicles
// @Tags         vehicle
// @Accept       json
// @Produce      json
// @Param        page    query    int  false  "Page number"
// @Param        page_size    query    int  false  "Page size"
// @Param        sort    query    int  false  "Sort field example '-name' or 'name'"
// @Param        mark    query    string  false  "Filter mark"
// @Param        mark_lk    query    string  false  "Filter mark like"
// @Param        mark    query    string  false  "Filter mark"
// @Param        mark_lk    query    string  false  "Filter mark like"
// @Param        model    query    string  false  "Filter model"
// @Param        model_lk    query    string  false  "Filter model like"
// @Param        regNum    query    string  false  "Filter regNum"
// @Param        regNum_lk    query    string  false  "Filter regNum like"
// @Param        year    query    int  false  "Filter mark"
// @Param        year_gt    query    int  false  "Filter year grater than"
// @Param        year_lt    query    int  false  "Filter year less than"
// @Param        year_btw    query    string  false  "Filter year between two years example '2000:2005'"
// @Param        owner_name    query    string  false  "Filter owner_name"
// @Param        owner_name_lk    query    string  false  "Filter owner_name like"
// @Param        owner_surname    query    string  false  "Filter owner_surname"
// @Param        owner_surname_lk    query    string  false  "Filter owner_surname like"
// @Param        owner_patronymic    query    string  false  "Filter owner_patronymic"
// @Param        owner_patronymic_lk    query    string  false  "Filter owner_patronymic like"
// @Success      200  {array}   pagination.PaginatedResponse[[]golang_vehicles.Vehicle]
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /vehicles [get]
func (h *Handler) GetAllVehicles(c *gin.Context) {
	qs := c.Request.URL.Query()

	validator := validators.NewValidator()
	pagination := pagination.NewPagination(qs, validator)

	filter := filter.NewFilter()

	var emptyValidates []validators.ValidateFunc[string]
	regNumValidates := []validators.ValidateFunc[string]{validators.ValidateRegNum(validator)}
	yearValidates := []validators.ValidateFunc[int]{validators.ValidateNumber(validator, 0, 1900)}

	filter.HandleAddString(qs, emptyValidates, "mark", "mark")
	filter.HandleAddString(qs, emptyValidates, "model", "model")
	filter.HandleAddString(qs, regNumValidates, "regNum", "reg_num")
	filter.HandleAddInt(qs, validator, yearValidates, "year", "year")
	filter.HandleAddString(qs, emptyValidates, "owner_name", "owner.name")
	filter.HandleAddString(qs, emptyValidates, "owner_surname", "owner.surname")
	filter.HandleAddString(qs, emptyValidates, "owner_patronymic", "owner.patronymic")

	if !validator.IsValid() {
		c.JSON(http.StatusBadRequest, validator.GetErrorsDict())
		return
	}

	PaginatedResponse, err := h.services.Vehicle.GetAllVehicles(filter, pagination)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse)
}

// @Summary      Delete vehicles
// @Tags         vehicle
// @Success      204
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /vehicles/{id} [delete]
func (h *Handler) DeleteVehicle(c *gin.Context) {
	vehicleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.DeleteVehicle(vehicleId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
