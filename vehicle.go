package golang_vehicles

type Vehicle struct {
	VehicleId int    `json:"id" db:"id"`
	RegNum    string `json:"regNum" db:"reg_num" binding:"required"`
	Mark      string `json:"mark" db:"mark" binding:"required"`
	Model     string `json:"model" db:"model" binding:"required"`
	Year      int    `json:"year" db:"year"`
	Owner     Owner  `json:"owner"`
}

type Owner struct {
	OwnerId    int    `json:"id" db:"id" `
	Name       string `json:"name" db:"name" `
	Surname    string `json:"surname" db:"surname" `
	Patronymic string `json:"patronymic" db:"patronymic"`
}

type UpdateVehicle struct {
	RegNum string      `json:"regNum" db:"reg_num"`
	Mark   string      `json:"mark" db:"mark"`
	Model  string      `json:"model" db:"model"`
	Year   int         `json:"year" db:"year"`
	Owner  UpdateOwner `json:"owner"`
}

type UpdateOwner struct {
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}
