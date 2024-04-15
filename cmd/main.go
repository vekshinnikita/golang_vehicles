package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vekshinnikita/golang_vehicles"
	"github.com/vekshinnikita/golang_vehicles/pkg/handler"
	"github.com/vekshinnikita/golang_vehicles/pkg/repository"
	"github.com/vekshinnikita/golang_vehicles/pkg/service"
)

//	@title		Golang Vehicle
//	@version	1.0

//	@host		localhost:8000
//	@BasePath	/api

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Info("Starting server...")

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	dbConfig := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	// command := fmt.Sprintf("-path ./schema -database 'postgres://%s:%s@%s:%s/%s?sslmode=%s' up",
	// 	dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	// )

	// cmd := exec.Command("migrate", "-path", "./schema", "-database", dbUrl, "up")
	// cmd := exec.Command("migrate", command)
	// out, err := cmd.Output()

	// fmt.Println(string(out))

	// if err != nil {
	// 	logrus.Fatalf("error ocurred while migrating db: %s", err.Error())
	// }

	db, err := repository.NewPostgresDB(dbConfig)

	if err != nil {
		logrus.Fatalf("error ocurred while initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(golang_vehicles.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error ocurred while running http server %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
